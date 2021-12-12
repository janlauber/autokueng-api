package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/janlauber/autokueng-api/controllers"
	"github.com/janlauber/autokueng-api/database"
	"github.com/janlauber/autokueng-api/models"
	"github.com/janlauber/autokueng-api/routes"
	"github.com/janlauber/autokueng-api/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	username string
	password string
	host     string
	dbName   string
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
	database.DBConn, err = gorm.Open("postgres", dbUri)
	if err != nil {
		panic(err)
	}

	database.DBConn.Debug().AutoMigrate(&models.User{})
	database.DBConn.Debug().AutoMigrate(&models.News{})

	util.InfoLogger.Println("Database connection initialized to: " + dbName)
}

func initEnvs() {
	controllers.SecretKey = os.Getenv("JWT_SECRET_KEY")
	if controllers.SecretKey == "" {
		util.ErrorLogger.Println("JWT_SECRET_KEY is not set")
		panic("stopping application...")
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
		controllers.UserAdmin = true
	}

}

func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("AutoKueng API")
	})

	routes.Setup(app)

	app.Listen(":8000")

	log.Printf("Server started on port 8000")

	defer database.DBConn.Close()
}
