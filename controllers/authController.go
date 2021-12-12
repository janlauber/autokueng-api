package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/janlauber/autokueng-api/database"
	"github.com/janlauber/autokueng-api/models"
	"github.com/janlauber/autokueng-api/util"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey string
var UserAdmin bool

func Register(c *fiber.Ctx) error {

	if UserAdmin {
		var data map[string]string

		if err := c.BodyParser(&data); err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
			return err
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)

		user := models.User{
			Username: data["username"],
			Password: string(password),
		}

		database.DBConn.Create(&user)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User created successfully",
		})
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User registration is disabled",
		})
	}
}

func ResetPassword(c *fiber.Ctx) error {

	if UserAdmin {

		var data map[string]string

		if err := c.BodyParser(&data); err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
			return err
		}

		var user models.User

		database.DBConn.Where("username = ?", data["username"]).First(&user)

		if user.Id == 0 {
			c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
			return errors.New("User not found")
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)

		user.Password = string(password)

		database.DBConn.Save(&user)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Password reset successfully",
		})
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User password reset is disabled",
		})
	}
}

func Login(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
		return err
	}

	user := models.User{}
	database.DBConn.Where("username = ?", data["username"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
		return nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid password",
		})
		return nil
	}

	claims := jwt.MapClaims{
		"id":   user.Id,
		"name": user.Username,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
		return err
	}

	cookie := fiber.Cookie{
		Name:   "jwt-autokueng-api",
		Value:  tokenString,
		MaxAge: 86400,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":       user.Id,
		"username": user.Username,
	})

}

func User(c *fiber.Ctx) error {
	if err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "authorized",
	})
}

func Logout(c *fiber.Ctx) error {

	if err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged out successfully",
	})
}

func CheckAuth(c *fiber.Ctx) error {
	// validate JWT token and authorize user
	cookie := c.Cookies("jwt-autokueng-api")

	if cookie == "" {
		util.WarningLogger.Println("IP " + c.IP() + " tried to access without a invalid token")
		return errors.New("no cookie")
	}

	// Parse the token and validate it
	token, _ := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	claims := token.Claims.(jwt.MapClaims)

	var user models.User

	database.DBConn.Where("id = ?", claims["id"]).First(&user)

	if user.Id == 0 {
		util.WarningLogger.Println("IP " + c.IP() + " tried to access with an invalid user")
		return errors.New("user not found")
	}

	return nil
}
