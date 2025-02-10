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
	if result := database.DbConn.Select(clause.Associations).Delete(&model.Bill{}, "id =?", id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No Bill found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No Bill found")
		}
		log.Printf("Error deleting Bill %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting Bill")
	}

	return c.Status(200).JSON("Bill deleted successfully")
}
