package controller

import (
	"errors"
	"log"
	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetDomain(c *fiber.Ctx) error {

	userId := c.Params("userId")

	log.Print(userId)
	log.Print(c.Get("userId"))

	if err := database.DbConn.First(&model.User{
		ID: userId,
	}).Error; err != nil {
		log.Print(err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(403).JSON(fiber.Map{
				"error": "Unauthenticated, User not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error" + err.Error(),
		})
	}

	var domain model.Domain
	if err := database.DbConn.Find(&domain, "user_id=?", userId).Error; err != nil {
		return c.Status(500).JSONP(fiber.Map{
			"error": "Internal server error" + err.Error(),
		})
	}
	log.Print(domain)
	return c.Status(200).JSON(domain)
}
