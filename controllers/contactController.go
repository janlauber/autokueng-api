package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	CaptchaSecret string
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

	var body ContactRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := CheckRecaptcha(body.RecaptchaResponse); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": "true",
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
