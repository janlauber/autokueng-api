package controllers

import (
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/janlauber/autokueng-api/database"
	"github.com/janlauber/autokueng-api/models"
)

func GetServices(c *fiber.Ctx) error {
	db := database.DBConn
	var services []models.Service
	db.Find(&services)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}

	// sort services by id
	sort.Slice(services, func(i, j int) bool {
		return services[i].ID < services[j].ID
	})
	return c.Status(200).JSON(services)
}

func GetService(c *fiber.Ctx) error {
	// get the id from the url
	id := c.Params("id")

	db := database.DBConn
	var service models.Service
	db.First(&service, id)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	if service.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Service not found"})
	}
	return c.Status(200).JSON(service)
}

func CreateService(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	db := database.DBConn

	// parse body to service model
	var service models.Service
	if err := c.BodyParser(&service); err != nil {
		return c.Status(400).JSON(err)
	}

	// create the service
	db.Create(&service)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(201).JSON(service)
}

func UpdateService(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	id := c.Params("id")
	db := database.DBConn

	// parse body to service model
	var service models.Service
	if err := c.BodyParser(&service); err != nil {
		return c.Status(400).JSON(err)
	}

	// check if service exists
	var count int64
	db.Model(&models.Service{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Service not found"})
	}

	// update the service
	db.Model(&models.Service{}).Where("id = ?", id).Updates(service)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}

	// find the updated service
	db.First(&service, id)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}

	return c.Status(200).JSON(service)
}

func DeleteService(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	id := c.Params("id")
	db := database.DBConn

	// check if service exists
	var count int64
	db.Model(&models.Service{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Service not found"})
	}

	// delete the service
	db.Where("id = ?", id).Delete(&models.Service{})
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(204).JSON(fiber.Map{})
}
