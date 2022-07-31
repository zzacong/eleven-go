package models

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/zzacong/eleven-go/go-fiber-crm/database"
	"gorm.io/gorm"
)

type Lead struct {
	gorm.Model
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

func GetLeads(c *fiber.Ctx) error {
	db := database.DB
	var leads []Lead
	db.Find(&leads)
	return c.JSON(leads)
}

func GetLead(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var lead Lead
	result := db.First(&lead, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fiber.ErrNotFound
	}
	return c.JSON(lead)
}

func NewLead(c *fiber.Ctx) error {
	db := database.DB
	lead := new(Lead)
	if err := c.BodyParser(lead); err != nil {
		return fiber.ErrBadRequest
	}
	db.Create(&lead)
	return c.Status(201).JSON(lead)
}

func DeleteLead(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var lead Lead
	result := db.First(&lead, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fiber.ErrNotFound
	}
	db.Delete(&lead)
	return c.SendString("Lead deleted successfully")
}
