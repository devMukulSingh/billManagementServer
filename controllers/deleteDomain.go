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

	if result := database.DbConn.Delete(&model.Domain{},"id =?",id); result.Error != nil {
		if errors.Is(result.Error,gorm.ErrRecordNotFound){
			log.Printf("No domain found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No domain found")
		}
		log.Printf("Error deleting domain %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting domain")
	}

	return c.Status(200).JSON("domain deleted successfully")
}
