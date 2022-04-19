package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/janlauber/autokueng-api/controllers"
	"github.com/janlauber/autokueng-api/database"
	"github.com/janlauber/autokueng-api/models"
	"github.com/janlauber/autokueng-api/routes"
	"github.com/janlauber/autokueng-api/util"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	username   string
	password   string
	host       string
	dbName     string
	corsString string
)

func init() {
	util.InitLoggers()
	initEnvs()
	initDB()
}

func initDB() {
	// Initialize DD
	var err error

	// build connection string
	dbUri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, username, password, dbName)

	// connect to database
	database.DBConn, err = gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	database.DBConn.Debug().AutoMigrate(&models.User{})
	database.DBConn.Debug().AutoMigrate(&models.News{})
	database.DBConn.Debug().AutoMigrate(&models.Service{})
	database.DBConn.Debug().AutoMigrate(&models.Member{})
	database.DBConn.Debug().AutoMigrate(&models.GalleryImage{})
	database.DBConn.Debug().AutoMigrate(&models.Link{})

	util.InfoLogger.Println("Database connection initialized to: " + dbName)
}

func initEnvs() {
	controllers.SecretKey = os.Getenv("JWT_SECRET_KEY")
	if controllers.SecretKey == "" {
		util.WarningLogger.Println("JWT_SECRET_KEY is not set")
		controllers.SecretKey = util.GenerateRandomString(32)
		util.InfoLogger.Println("JWT_SECRET_KEY is set to default: " + controllers.SecretKey)
	}

	corsString = os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsString == "" {
		util.WarningLogger.Println("CORS_ALLOWED_ORIGINS is not set")
		corsString = "*"
		util.InfoLogger.Println("CORS_ALLOWED_ORIGINS is set to default: " + corsString)
	}
	username = os.Getenv("DB_USERNAME")
	if username == "" {
		util.ErrorLogger.Println("DB_USERNAME is not set")
		panic("stopping application...")
	}
	// could be empty
	password = os.Getenv("DB_PASSWORD")
	if password == "" {
		util.WarningLogger.Println("DB_PASSWORD is not set")
	}

	host = os.Getenv("DB_HOST")
	if host == "" {
		util.ErrorLogger.Println("DB_HOST is not set")
		panic("stopping application...")
	}
	dbName = os.Getenv("DB_NAME")
	if dbName == "" {
		util.ErrorLogger.Println("DB_NAME is not set")
		panic("stopping application...")
	}
	if os.Getenv("USER_ADMIN") == "enabled" {
		util.InfoLogger.Println("User administration is enabled")
		controllers.UserAdmin = true
	}
	controllers.CaptchaSecret = os.Getenv("CAPTCHA_SECRET")
	if controllers.CaptchaSecret == "" {
		util.ErrorLogger.Println("CAPTCHA_SECRET is not set")
		panic("stopping application...")
	}

}

func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     corsString,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("AutoKueng API")
	})

	routes.Setup(app)

	app.Listen(":8000")
}
