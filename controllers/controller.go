package controller

import "github.com/gofiber/fiber/v2"

func PostBillController(c *fiber.Ctx) error{
	return c.JSON("Hello worlld")
}