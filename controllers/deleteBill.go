package controller

import (
	"errors"
	"log"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func DeleteBill(c *fiber.Ctx) error {
	id := c.Params("id")
	var existingBill model.Bill

	if result := database.DbConn.First(&existingBill, "id=?", id); result != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No bill found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No Record found")
		}
	}

	if result := database.DbConn.Select(clause.Associations).Delete(&existingBill); result.Error != nil {
		log.Printf("Error deleting Bill %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting Bill")
	}

	return c.Status(200).JSON("Bill deleted successfully")
}
