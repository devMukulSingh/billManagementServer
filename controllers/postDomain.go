package controller

import (
	"errors"
	"log"
	"gorm.io/gorm"
	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
)

type Domain struct {
	Domain string `json:"domain"`
}

func PostDomain(c *fiber.Ctx) error {

	body := new(Domain)
	userId := c.Params("userId")

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}

	result := database.DbConn.Create(&model.Domain{
		Name: body.Domain,
		UserID: userId,
	})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			log.Print(result.Error.Error())
			return c.Status(409).JSON(fiber.Map{
				"error":"Domain already exists, try another",
			})
		}
		log.Printf("Error in saving Domain into db %s", result.Error.Error())
		return c.Status(500).JSON("Internal server error")
	}

	return c.Status(201).JSON(fiber.Map{
		"msg": "domain created successfully",
	})
}
