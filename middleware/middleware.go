package middleware

import (
	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/gofiber/fiber/v2"
)

func ValidateUser(c *fiber.Ctx) error {
	var params types.Param

	if err := c.ParamsParser(&params); err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error": "Invalid userId",
		})
	}

	return c.Next();
}
