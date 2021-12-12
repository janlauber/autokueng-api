package controllers

import "github.com/gofiber/fiber/v2"

func GetServices(c *fiber.Ctx) error {
	return c.Status(200).JSON(map[string]interface{}{
		"message": "Get Services",
	})
}

func GetService(c *fiber.Ctx) error {
	return c.Status(200).JSON(map[string]interface{}{
		"message": "Get Service",
	})
}

func CreateService(c *fiber.Ctx) error {
	return c.Status(200).JSON(map[string]interface{}{
		"message": "Create Services",
	})
}

func UpdateService(c *fiber.Ctx) error {
	return c.Status(200).JSON(map[string]interface{}{
		"message": "Update Service",
	})
}

func DeleteService(c *fiber.Ctx) error {
	return c.Status(200).JSON(map[string]interface{}{
		"message": "Delete Services",
	})
}
