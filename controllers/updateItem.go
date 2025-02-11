package controller

import (
	"log"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
)

func UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")

	var existingItem model.Item;

	if result := database.DbConn.First(&existingItem, "id = ?", id); result.Error != nil {
		log.Printf("No item found %s", result.Error.Error())
		return c.Status(400).JSON("No item found")
	}

	body := new(Item)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}
	if result := database.DbConn.Model(&existingItem).Updates(
		model.Item{
			Name:     body.Name,
			Rate:     body.Rate,
			Amount:   body.Amount,
			Quantity: body.Rate,
		},
	); result.Error != nil {
		log.Printf("Error updating item %s", result.Error.Error())
		return c.Status(500).JSON("Error updating item")
	}

	return c.Status(200).JSON("item updated successfully")
}
