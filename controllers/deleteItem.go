package controller

import (
	"errors"
	"log"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DeleteItem(c *fiber.Ctx) error {
	id := c.Params("id")

	var existingItem model.Item

	if result := database.DbConn.First(&existingItem, "id=?", id); result != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No item found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No Record found")
		}
	}

	if result := database.DbConn.Delete(&existingItem); result.Error != nil {
		log.Printf("Error deleting item %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting item")
	}

	return c.Status(200).JSON("item deleted successfully")
}
