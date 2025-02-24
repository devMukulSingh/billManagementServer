package controller

import (
	"errors"
	"log"

	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllDomains(c *fiber.Ctx) error {
	userId := c.Params("userId")
	var domains []model.Domain
	if err := database.DbConn.Find(&domains, "user_id=?", userId).Error; err != nil {
		return c.Status(500).JSONP(fiber.Map{
			"error": "Internal server error " + err.Error(),
		})
	}
	return c.Status(200).JSON(domains)
}

func GetDomain(c *fiber.Ctx) error {

	userId := c.Params("userId")
	domainId := c.Params("domainId")

	var domain model.Domain
	if err := database.DbConn.Where("id = ? AND user_id = ?",domainId,userId).Limit(1).Find(&domain).Error; err != nil {
		return c.Status(500).JSONP(fiber.Map{
			"error": "Internal server error " + err.Error(),
		})
	}
	return c.Status(200).JSON(domain)
}


func PostDomain(c *fiber.Ctx) error {

	body := new(types.Domain)
	userId := c.Params("userId")

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}

	result := database.DbConn.Create(&model.Domain{
		Name:   body.DomainName,
		UserID: userId,
	})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			log.Print(result.Error.Error())
			return c.Status(409).JSON(fiber.Map{
				"error": "Domain already exists, try another",
			})
		}
		log.Printf("Error in saving Domain into db %s", result.Error.Error())
		return c.Status(500).JSON("Internal server error")
	}

	return c.Status(201).JSON(fiber.Map{
		"msg": "domain created successfully",
	})
}

func UpdateDomain(c *fiber.Ctx) error {
	domainId := c.Params("domainId")

	var existingDomain model.Domain

	if result := database.DbConn.First(&existingDomain, "id = ?", domainId); result.Error != nil {
		log.Printf("No domain found %s", result.Error.Error())
		return c.Status(400).JSON("No domain found")
	}

	body := new(types.Domain)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}
	if result := database.DbConn.Model(&existingDomain).Update("name", body.DomainName); result.Error != nil {
		log.Printf("Error updating domain %s", result.Error.Error())
		return c.Status(500).JSON("Error updating domain")
	}

	return c.Status(200).JSON("domain updated successfully")
}

func DeleteDomain(c *fiber.Ctx) error {
	domainId := c.Params("domainId")
	var existingDomain model.Domain

	if result := database.DbConn.First(&existingDomain, "id=?", domainId); result != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No domain found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No record found")
		}
	}

	if result := database.DbConn.Delete(&existingDomain); result.Error != nil {

		log.Printf("Error deleting domain %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting domain")
	}

	return c.Status(200).JSON("domain deleted successfully")
}
