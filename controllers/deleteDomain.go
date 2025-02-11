package controller

import (
	"errors"
	"log"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DeleteDomain(c *fiber.Ctx) error {
	id := c.Params("id")
	var existingDomain model.Domain

	if result := database.DbConn.First(&existingDomain, "id=?", id); result != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No domain found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No record found")
		}
	}

	if result := database.DbConn.Delete(&existingDomain); result.Error != nil {

		log.Printf("Error deleting domain %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting domain")
	}

	return c.Status(200).JSON("domain deleted successfully")
}
