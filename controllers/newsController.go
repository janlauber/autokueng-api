package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/janlauber/autokueng-api/database"
	"github.com/janlauber/autokueng-api/models"
)

func GetNews(c *fiber.Ctx) error {
	db := database.DBConn
	var news []models.News
	db.Find(&news)
	return c.Status(200).JSON(news)
}

func CreateNews(c *fiber.Ctx) error {

	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// create the news
	db := database.DBConn
	var count int64
	db.Model(&models.News{}).Count(&count)
	// if no news is found in the database, create a new one
	if count == 0 {
		news := new(models.News)
		if err := c.BodyParser(news); err != nil {
			return c.Status(400).JSON(err)
		}
		db.Create(&news)
		return c.Status(201).JSON(news)
	} else {
		return c.Status(400).JSON(fiber.Map{"error": "News already exists"})
	}
}

func UpdateNews(c *fiber.Ctx) error {

	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	db := database.DBConn
	news := new(models.News)
	db.First(&news)
	if news.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "News not found"})
	}
	if err := c.BodyParser(news); err != nil {
		return c.Status(400).JSON(err)
	}
	db.Save(&news)
	return c.Status(200).JSON(news)
}

func DeleteAllNews(c *fiber.Ctx) error {

	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	db := database.DBConn
	db.Delete(&models.News{})
	return c.Status(200).JSON(fiber.Map{"message": "News deleted"})
}
