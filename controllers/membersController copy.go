package controllers

import (
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/janlauber/autokueng-api/database"
	"github.com/janlauber/autokueng-api/models"
)

func GetMembers(c *fiber.Ctx) error {
	db := database.DBConn
	var members []models.Member
	db.Find(&members)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}

	// sort members by id
	sort.Slice(members, func(i, j int) bool {
		return members[i].ID < members[j].ID
	})
	return c.Status(200).JSON(members)
}

func GetMember(c *fiber.Ctx) error {
	// get the id from the url
	id := c.Params("id")

	db := database.DBConn
	var member models.Member
	db.First(&member, id)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	if member.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Member not found"})
	}
	return c.Status(200).JSON(member)
}

func CreateMember(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	db := database.DBConn

	// parse body to member model
	var member models.Member
	if err := c.BodyParser(&member); err != nil {
		return c.Status(400).JSON(err)
	}

	// create the member
	db.Create(&member)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(201).JSON(member)
}

func UpdateMember(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	id := c.Params("id")
	db := database.DBConn

	// parse body to member model
	var member models.Member
	if err := c.BodyParser(&member); err != nil {
		return c.Status(400).JSON(err)
	}

	// check if member exists
	var count int64
	db.Model(&models.Member{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Member not found"})
	}

	// update the member
	db.Model(&models.Member{}).Where("id = ?", id).Updates(member)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}

	// return the updated member
	db.First(&member, id)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(200).JSON(member)
}

func DeleteMember(c *fiber.Ctx) error {
	if _, err := CheckAuth(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	id := c.Params("id")
	db := database.DBConn

	// check if member exists
	var count int64
	db.Model(&models.Member{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Member not found"})
	}

	// delete the member
	db.Delete(&models.Member{}, id)
	if db.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": db.Error})
	}
	return c.Status(204).JSON(fiber.Map{})
}
