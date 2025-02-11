package controller

import (
	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
	"log"
)


func PostItem(c *fiber.Ctx) error {
	body := new(Item)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}

	result := database.DbConn.Create(&model.Item{
		Name: body.Name,
		Rate: body.Rate,
		Amount: body.Amount,
		Quantity: body.Quantity,
	})
	if result.Error != nil {
		log.Printf("Error in saving Items into db %s", result.Error.Error())
		return c.Status(500).JSON("Internal server error")
	}

	return c.Status(201).JSON(fiber.Map{
		"msg": "Item created successfully",
	})
}
