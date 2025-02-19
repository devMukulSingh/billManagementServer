package controller

import (
	"errors"
	"log"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Distributor struct {
	Distributor string `json:"distributor"`
	DomainID    string `json:"domainId"`
}

func PostDistributor(c *fiber.Ctx) error {

	body := new(Distributor)
	userId := c.Params("userId")

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing body %s", err.Error())
		return c.Status(400).JSON("Error parssing body")
	}

	result := database.DbConn.Create(&model.Distributor{
		Name:     body.Distributor,
		DomainID: body.DomainID,
		UserID: userId,
	})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			log.Print(result.Error.Error())
			return c.Status(409).JSON(fiber.Map{
				"error": "Distributor already exists, try another",
			})
		}
		log.Printf("Error in saving Distributor into db %s", result.Error.Error())
		return c.Status(500).JSON("Internal server error")
	}
	return c.Status(201).JSON(fiber.Map{
		"msg": "distributor created successfully",
	})

}
