package controllers

import (
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/janlauber/autokueng-api/database"
	"github.com/janlauber/autokueng-api/models"
)

func GetLinks(c *fiber.Ctx) error {
	db := database.DBConn
	var links []models.Link
	db.Find(&links)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}

	// sort links by id
	sort.Slice(links, func(i, j int) bool {
		return links[i].ID < links[j].ID
	})
	return c.Status(200).JSON(links)
}

func GetLink(c *fiber.Ctx) error {
	// get the id from the url
	id := c.Params("id")

	db := database.DBConn
	var link models.Link
	db.First(&link, id)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	if link.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Link not found"})
	}
	return c.Status(200).JSON(link)
}

func CreateLink(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	db := database.DBConn

	// parse body to link model
	var link models.Link
	if err := c.BodyParser(&link); err != nil {
		return c.Status(400).JSON(err)
	}

	// create the link
	db.Create(&link)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(201).JSON(link)
}

func UpdateLink(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	id := c.Params("id")
	db := database.DBConn

	// parse body to link model
	var link models.Link
	if err := c.BodyParser(&link); err != nil {
		return c.Status(400).JSON(err)
	}

	// check if link exists
	var count int64
	db.Model(&models.Link{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Link not found"})
	}

	// update the link
	db.First(&link, id).Update(c.BodyParser(&link)).Save(&link)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(200).JSON(link)
}

func DeleteLink(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	id := c.Params("id")
	db := database.DBConn

	// check if link exists
	var count int64
	db.Model(&models.Link{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Link not found"})
	}

	// delete the link
	db.Where("id = ?", id).Delete(&models.Link{})
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(204).JSON(fiber.Map{})
}
