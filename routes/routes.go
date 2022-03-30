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

	// Members
	v1.Get("/members", controllers.GetMembers)
	v1.Get("/members/:id", controllers.GetMember)
	v1.Post("/members", controllers.CreateMember)
	v1.Put("/members/:id", controllers.UpdateMember)
	v1.Delete("/members/:id", controllers.DeleteMember)

	// GalleryImages
	v1.Get("/gallery-images", controllers.GetGalleryImages)
	v1.Get("/gallery-images/:id", controllers.GetGalleryImage)
	v1.Post("/gallery-images", controllers.CreateGalleryImage)
	v1.Put("/gallery-images/:id", controllers.UpdateGalleryImage)
	v1.Delete("/gallery-images/:id", controllers.DeleteGalleryImage)

	// Links
	v1.Get("/links", controllers.GetLinks)
	v1.Get("/links/:id", controllers.GetLink)
	v1.Post("/links", controllers.CreateLink)
	v1.Put("/links/:id", controllers.UpdateLink)
	v1.Delete("/links/:id", controllers.DeleteLink)
}
