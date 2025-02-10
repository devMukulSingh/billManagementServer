package controller

import (
	"log"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
)

func UpdateDistributor(c *fiber.Ctx) error {
	id := c.Params("id")

	var exisitingDistributor model.Distributor

	if result := database.DbConn.First(&exisitingDistributor, "id = ?", id); result.Error != nil {
		log.Printf("No distributor found %s", result.Error.Error())
		return c.Status(400).JSON("No distributor found")
	}

	body  := new(Distributor)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}
	if result := database.DbConn.Model(&exisitingDistributor).Updates(
		model.Distributor{
			Name: body.Distributor,
			DomainID: body.DomainID,
		},
	); result.Error != nil {
		log.Printf("Error updating distributor %s", result.Error.Error())
		return c.Status(500).JSON("Error updating distributor")
	}

	return c.Status(200).JSON("distributor updated successfully")
}
