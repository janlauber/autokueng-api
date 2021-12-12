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
	// TODO: if middleware CookieAuthRequired is nil then go to handler
	v1.Get("/user", controllers.User)

	// News
	v1.Get("/news", controllers.GetNews)
	v1.Put("/news", controllers.CreateNews)
	v1.Post("/news", controllers.UpdateNews)
	v1.Delete("/news", controllers.DeleteAllNews)
}
