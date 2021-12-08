package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/janlauber/autokueng-api/controllers"
)

func Setup(app *fiber.App) {
	app.Post("/api/v1/register", controllers.Register)
	app.Post("/api/v1/login", controllers.Login)
	app.Get("/api/v1/user", controllers.User)
	app.Post("/api/v1/logout", controllers.Logout)

	app.Get("/api/v1/news", controllers.GetNews)
	app.Put("/api/v1/news", controllers.CreateNews)
	app.Post("/api/v1/news", controllers.UpdateNews)
	app.Delete("/api/v1/news", controllers.DeleteAllNews)
}
