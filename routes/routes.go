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
}
