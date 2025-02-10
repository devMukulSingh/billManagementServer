package controller

import (
	"log"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
)

func UpdateDomain(c *fiber.Ctx) error {
	id := c.Params("id")

	var existingDomain model.Domain

	if result := database.DbConn.First(&existingDomain, "id = ?", id); result.Error != nil {
		log.Printf("No domain found %s", result.Error.Error())
		return c.Status(400).JSON("No domain found")
	}

	body  := new(Domain)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}
	if result := database.DbConn.Model(&existingDomain).Update("name", body.Domain); result.Error != nil {
		log.Printf("Error updating domain %s", result.Error.Error())
		return c.Status(500).JSON("Error updating domain")
	}

	return c.Status(200).JSON("domain updated successfully")
}
