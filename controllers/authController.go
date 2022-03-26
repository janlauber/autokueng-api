package controllers

import (
	"errors"
	"strings"
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

		if len(data["username"]) == 0 {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username is required",
			})
			return nil
		}

		if len(data["password"]) == 0 {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Password is required",
			})
			return nil
		}

		password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
			return err
		}

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
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":       user.Id,
		"username": user.Username,
		"token":    tokenString,
	})

}

func Auth(c *fiber.Ctx) error {
	uid, err := CheckAuth(c)
	if err != nil {
		util.ErrorLogger.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var user models.User

	database.DBConn.Where("id = ?", uid).First(&user)
	if user.Id == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":       user.Id,
		"username": user.Username,
	})
}

func CheckAuth(c *fiber.Ctx) (int64, error) {
	// validate JWT token from Baerer authorization header
	var token *jwt.Token
	var tokenString string
	var err error

	bearerToken := c.Get("Authorization")

	bearerTokenSplit := strings.Split(bearerToken, " ")
	if len(bearerTokenSplit) != 2 {
		return 0, errors.New("Invalid bearer token")
	} else {
		tokenString = bearerTokenSplit[1]
	}

	if len(tokenString) == 0 {
		return 0, errors.New("Invalid bearer token")
	}

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid token")
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return 0, errors.New("Invalid token")
	}

	if token == nil || !token.Valid {
		return 0, errors.New("Invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["exp"] == nil {
		return 0, nil
	} else {
		exp := claims["exp"]
		// convert exp to int64
		expInt64 := int64(exp.(float64))
		if expInt64 < time.Now().Unix() {
			return 0, nil
		}
	}

	var user models.User

	database.DBConn.Where("id = ?", claims["id"]).First(&user)

	if user.Id == 0 {
		util.WarningLogger.Println("IP " + c.IP() + " tried to access with an invalid user")
		return 0, errors.New("user not found")
	}

	return user.Id, nil
}
