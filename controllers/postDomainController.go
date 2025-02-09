package controller
import (
	"log"
	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
)

func PostDomainController(c *fiber.Ctx) error {

	type Domain struct {
		Domain string `json:"domain"`
	}
	req := new(Domain)

	if err := c.BodyParser(req); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}

	result := database.DbConn.Create(&model.Domain{
		Name: req.Domain,
	})
	if result.Error != nil {
		log.Printf("Error in saving Domain into db %s", result.Error.Error())
		return c.Status(500).JSON("Internal server error")	
	}

	return c.Status(201).JSON(fiber.Map{
		"msg":"domain created successfully",
	})
}
