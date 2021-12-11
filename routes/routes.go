package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/janlauber/autokueng-api/controllers"
)

func Setup(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Users
	v1.Post("/register", controllers.Register)
	v1.Post("/login", controllers.Login)
	v1.Post("/logout", controllers.Logout)
	auth := v1.Group("/auth", controllers.CookieAuthRequired())
	auth.Get("/user", controllers.User)

	app.Get("/api/v1/news", controllers.GetNews)
	app.Put("/api/v1/news", controllers.CreateNews)
	app.Post("/api/v1/news", controllers.UpdateNews)
	app.Delete("/api/v1/news", controllers.DeleteAllNews)
}
