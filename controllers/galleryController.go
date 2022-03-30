package controllers

import (
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/janlauber/autokueng-api/database"
	"github.com/janlauber/autokueng-api/models"
)

func GetGalleryImages(c *fiber.Ctx) error {
	db := database.DBConn
	var galleryImages []models.GalleryImage
	db.Find(&galleryImages)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}

	// sort galleryImage by id
	sort.Slice(galleryImages, func(i, j int) bool {
		return galleryImages[i].ID < galleryImages[j].ID
	})
	return c.Status(200).JSON(galleryImages)
}

func GetGalleryImage(c *fiber.Ctx) error {
	// get the id from the url
	id := c.Params("id")

	db := database.DBConn
	var galleryImage models.GalleryImage
	db.First(&galleryImage, id)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	if galleryImage.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "GalleryImage not found"})
	}
	return c.Status(200).JSON(galleryImage)
}

func CreateGalleryImage(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	db := database.DBConn

	// parse body to GalleryImage model
	var galleryImage models.GalleryImage
	if err := c.BodyParser(&galleryImage); err != nil {
		return c.Status(400).JSON(err)
	}

	// create the galleryImage
	db.Create(&galleryImage)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(201).JSON(galleryImage)
}

func UpdateGalleryImage(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	id := c.Params("id")
	db := database.DBConn

	// parse body to gallerImage model
	var gallerImage models.GalleryImage
	if err := c.BodyParser(&gallerImage); err != nil {
		return c.Status(400).JSON(err)
	}

	// check if galleryImage exists
	var count int64
	db.Model(&models.GalleryImage{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Member not found"})
	}

	// update the galleryImage
	db.Model(&models.GalleryImage{}).Where("id = ?", id).Updates(gallerImage)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}

	// return the updated member
	db.First(&gallerImage, id)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(200).JSON(gallerImage)
}

func DeleteGalleryImage(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	id := c.Params("id")
	db := database.DBConn

	// check if galleryImage exists
	var count int64
	db.Model(&models.GalleryImage{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "GalleryImage not found"})
	}

	// delete the galleryImage
	db.Delete(&models.GalleryImage{}, id)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(204).JSON(fiber.Map{})
}
