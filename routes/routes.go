package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/janlauber/autokueng-api/controllers"
)

func Setup(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")
	admin := v1.Group("/admin")

	// Users
	v1.Post("/auth", controllers.Login)
	// TODO: if middleware CookieAuthRequired is nil then go to handler
	v1.Get("/auth", controllers.Auth)
	// Adminstuff (disabled)
	admin.Post("/register", controllers.Register)
	admin.Post("/reset-password", controllers.ResetPassword)

	// News
	v1.Get("/news", controllers.GetNews)
	v1.Post("/news", controllers.CreateNews)
	v1.Put("/news", controllers.UpdateNews)
	v1.Delete("/news", controllers.DeleteAllNews)

	// Services
	v1.Get("/services", controllers.GetServices)
	v1.Get("/services/:id", controllers.GetService)
	v1.Post("/services", controllers.CreateService)
	v1.Put("/services/:id", controllers.UpdateService)
	v1.Delete("/services/:id", controllers.DeleteService)

}
