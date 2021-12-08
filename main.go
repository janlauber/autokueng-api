package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/janlauber/autokueng-api/database"
	"github.com/janlauber/autokueng-api/models"
	"github.com/janlauber/autokueng-api/routes"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func initDB() {
	// Initialize DB

	var err error

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// build connection string
	dbUri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, username, password, dbName)

	// connect to database
	database.DBConn, err = gorm.Open("postgres", dbUri)
	if err != nil {
		panic(err)
	}

	database.DBConn.Debug().AutoMigrate(&models.User{})
	database.DBConn.Debug().AutoMigrate(&models.News{})

	log.Printf("Successfully connected to database %s\n", dbName)
}

func main() {

	initDB()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	routes.Setup(app)

	app.Listen(":8000")

	log.Printf("Server started on port 8000")

	defer database.DBConn.Close()
}
