package controller

import (
	"log"
	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
)

type Distributor struct {
	Distributor string `json:"distributor"`
	DomainID    string    `json:"domainId"`
}

func PostDistributor(c *fiber.Ctx) error {

	body := new(Distributor)
	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing body %s", err.Error())
		return c.Status(400).JSON("Error parssing body")
	}
	result := database.DbConn.Create(&model.Distributor{
		Name:     body.Distributor,
		DomainID: body.DomainID,
	}); 
	if result.Error != nil {
		log.Printf("Error in saving Distributor into db %s", result.Error.Error())
		return c.Status(500).JSON("Internal server error")
	}

	return c.Status(201).JSON(fiber.Map{
		"msg":"distributor created successfully",

	})
}
