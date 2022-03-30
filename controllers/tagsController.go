package controllers

import (
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/janlauber/autokueng-api/database"
	"github.com/janlauber/autokueng-api/models"
)

func GetTags(c *fiber.Ctx) error {
	db := database.DBConn
	var tags []models.Tag
	db.Find(&tags)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}

	// sort tags by id
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].ID < tags[j].ID
	})
	return c.Status(200).JSON(tags)
}

func GetTag(c *fiber.Ctx) error {
	// get the id from the url
	id := c.Params("id")

	db := database.DBConn
	var tag models.Tag
	db.First(&tag, id)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	if tag.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tag not found"})
	}
	return c.Status(200).JSON(tag)
}

func CreateTag(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	db := database.DBConn

	// parse body to tag model
	var tag models.Tag
	if err := c.BodyParser(&tag); err != nil {
		return c.Status(400).JSON(err)
	}

	// create the tag
	db.Create(&tag)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(201).JSON(tag)
}

func UpdateTag(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	id := c.Params("id")
	db := database.DBConn

	// parse body to tag model
	var tag models.Tag
	if err := c.BodyParser(&tag); err != nil {
		return c.Status(400).JSON(err)
	}

	// check if tag exists
	var count int64
	db.Model(&models.Tag{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Link not found"})
	}

	// update the tag
	db.First(&tag, id).Save(&tag)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(200).JSON(tag)
}

func DeleteTag(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	id := c.Params("id")
	db := database.DBConn

	// check if tag exists
	var count int64
	db.Model(&models.Tag{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tag not found"})
	}

	// delete the tag
	db.Where("id = ?", id).Delete(&models.Tag{})
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(204).JSON(fiber.Map{})
}
