package controller

import (
	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Domain struct {
	Domain string `json:"domain"`
}

func PostDomain(c *fiber.Ctx) error {
	body := new(Domain)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}

	result := database.DbConn.Create(&model.Domain{
		Name: body.Domain,
	})
	if result.Error != nil {
		log.Printf("Error in saving Domain into db %s", result.Error.Error())
		return c.Status(500).JSON("Internal server error")
	}

	return c.Status(201).JSON(fiber.Map{
		"msg": "domain created successfully",
	})
}
