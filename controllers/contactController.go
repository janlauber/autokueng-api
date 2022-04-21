package controllers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	CaptchaSecret string
	SMTP_Username string
	SMTP_Password string
	SMTP_Host     string
	SMTP_Port     string
	SMTP_From     string
	SMTP_To       string
	SMTP_SSL      bool
)

type ContactRequest struct {
	Firstname         string `json:"firstname"`
	Lastname          string `json:"lastname"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Subject           string `json:"subject"`
	Message           string `json:"message"`
	RecaptchaResponse string `json:"g-recaptcha-response"`
}

const siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

func SendContactform(c *fiber.Ctx) error {

	var reqbody ContactRequest
	if err := c.BodyParser(&reqbody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := CheckRecaptcha(reqbody.RecaptchaResponse); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Send email
	from := mail.Address{Name: "(Kontaktformular) " + reqbody.Firstname + " " + reqbody.Lastname, Address: reqbody.Email}
	to := mail.Address{Name: "", Address: SMTP_To}
	subject := reqbody.Subject
	body := fmt.Sprintf("Vorname: %s\nNachname: %s\nE-Mail: %s\nTelefon: %s\n\nNachricht: \n%s", reqbody.Firstname, reqbody.Lastname, reqbody.Email, reqbody.Phone, reqbody.Message)

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	servername := fmt.Sprintf("%s:%s", SMTP_Host, SMTP_Port)

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", SMTP_Username, SMTP_Password, host)

	if SMTP_SSL {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
		conn, err := tls.Dial("tcp", servername, tlsconfig)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		con, err := smtp.NewClient(conn, host)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := con.Auth(auth); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := con.Mail(from.Address); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := con.Rcpt(to.Address); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		wdata, err := con.Data()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		_, err = wdata.Write([]byte(message))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = wdata.Close()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		con.Quit()
	} else {
		// Connect to the server, authenticate, set the sender and recipient,
		// and send the email all in one step.
		if err := smtp.SendMail(servername, auth, from.Address, []string{to.Address}, []byte(message)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": "ok",
	})

}

func CheckRecaptcha(response string) error {
	req, err := http.NewRequest("POST", siteVerifyURL,
		strings.NewReader(fmt.Sprintf("secret=%s&response=%s", CaptchaSecret, response)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result["success"] != true {
		return fmt.Errorf("recaptcha failed")
	}

	return nil
}
