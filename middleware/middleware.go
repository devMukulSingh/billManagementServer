package middleware

import "github.com/gofiber/fiber/v2"

func ValidateUser(c *fiber.Ctx) error {
	param := struct {
		UserId string `params:"userId"`
	}{}

	if err := c.ParamsParser(&param); err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error": "Invalid userId",
		})
	}

	return c.Next();
}
