package controller

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/devMukulSingh/billManagementServer.git/valkeyCache"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllDomains(c *fiber.Ctx) error {

	userId := c.Params("userId")

	//valkey cache
	cache, err := valkeyCache.GetValue("domains:" + userId)
	if err != nil {
		if err.Error() != "valkey nil message" {
			log.Printf("Error in getting cached bills : %s", err)
		}
	} else {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString(cache)
	}

	var domains []model.Domain
	if err := database.DbConn.Find(&domains, "user_id=?", userId).Error; err != nil {
		return c.Status(500).JSONP(fiber.Map{
			"error": "Internal server error " + err.Error(),
		})
	}
	jsonDomain, err := json.Marshal(domains)
	if err != nil {
		log.Print("error converting to json")
	}
	if err := valkeyCache.SetValue("domains:"+userId, jsonDomain); err != nil {
		log.Printf("Error in setting value in valkey %s ", err.Error())
	}
	return c.Status(200).JSON(domains)
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

	//revalidate cache
	valkeyCache.Revalidate("domains:" + userId)
	return c.Status(201).JSON(fiber.Map{
		"msg": "domain created successfully",
	})
}

func UpdateDomain(c *fiber.Ctx) error {
	domainId := c.Params("domainId")
	userId := c.Params("userId")

	body := new(types.Domain)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}
	if result := database.DbConn.Model(&model.Domain{}).Where("id=? AND user_id",domainId,userId).Update("name", body.DomainName); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No domain found %s", result.Error.Error())
			return c.Status(400).JSON("No domain found")
		}
		log.Printf("Error updating domain %s", result.Error.Error())
		return c.Status(500).JSON("Error updating domain")
	}
	if err := valkeyCache.Revalidate("domains:" + userId); err != nil {
		log.Printf("Error in revalidating domains cache: %s", err)
	}
	return c.Status(200).JSON("domain updated successfully")
}

func DeleteDomain(c *fiber.Ctx) error {
	domainId := c.Params("domainId")
	userId := c.Params("userId")

	if err := database.DbConn.Where("id =?", domainId).Delete(&model.Domain{}).Error; err != nil {
		log.Printf("Error deleting domain %s", err.Error())

		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No domain found %s", err.Error())
			return c.Status(400).JSON("Error:No record found")
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return c.Status(405).JSON(fiber.Map{
				"error": "Delete associated bills and distributors to delete domain",
			})
		}

		return c.Status(500).JSON("Error deleting domain")
	}
	if err := valkeyCache.Revalidate("domains:" + userId); err != nil {
		log.Printf("Error in revalidating domains cache: %s", err)
	}
	return c.Status(200).JSON("domain deleted successfully")
}

func GetDomain(c *fiber.Ctx) error {

	userId := c.Params("userId")
	domainId := c.Params("domainId")

	var domain model.Domain
	if err := database.DbConn.Where("id = ? AND user_id = ?", domainId, userId).Limit(1).Find(&domain).Error; err != nil {
		return c.Status(500).JSONP(fiber.Map{
			"error": "Internal server error " + err.Error(),
		})
	}
	return c.Status(200).JSON(domain)
}
