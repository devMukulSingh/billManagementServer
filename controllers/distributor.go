package controller

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/devMukulSingh/billManagementServer.git/types"
	// "github.com/devMukulSingh/billManagementServer.git/valkeyCache"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func GetAllDistributors(c *fiber.Ctx) error {

	userId := c.Params("userId")
	page := c.Query("page")
	// cached, err := valkeyCache.GetValue("distributors:" + page + ":" + userId)
	
	// if err != nil {
	// 	if err.Error() != "valkey nil message" {
	// 		log.Printf("Error in getting cached bills : %s", err)
	// 	}
	// } else {
	// 	c.Set("Content-Type", "application/json")
	// 	return c.SendString(cached)
	// }

	type Distributor struct {
		ID        string          `json:"id" `
		CreatedAt time.Time       `json:"created_at"`
		Name      string          `json:"name" `
		Domain    json.RawMessage `json:"domain"`
	}

	var data []Distributor
	var count int64
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Print("Error converting page to Int")
	}
	if err := database.DbConn.Model(&model.Distributor{}).Count(&count).Offset((pageInt-1)*10).Limit(10).
		Joins("JOIN domains ON domains.id = distributors.domain_id").
		Select(`
		distributors.id,
		distributors.name,
		distributors.created_at,
		json_build_object(
		'id' , domains.id,
		'name',domains.name
		) as domain
	`).
		Where("distributors.user_id =?", userId).
		Scan(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error " + err.Error(),
		})
	}

	type Result struct {
		Data   []Distributor			`json:"data"`
		Count int64							`json:"count"`
	}
	result := Result{
		Data:   data,
		Count: count,
	}
	// jsonString, err := json.Marshal(result)
	// if err != nil {
	// 	log.Printf("Error converting into Json: %s", err)
	// 	return c.Status(200).JSON(result)
	// }
	// if err := valkeyCache.SetValue("distributors:"+page+":"+userId, jsonString); err != nil {
	// 	log.Printf("Error in setting value in valkey %s ", err.Error())
	// }
	return c.Status(200).JSON(result)

}

func PostDistributor(c *fiber.Ctx) error {

	body := new(types.Distributor)
	userId := c.Params("userId")

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing body %s", err.Error())
		return c.Status(400).JSON("Error parssing body")
	}

	result := database.DbConn.Create(&model.Distributor{
		DomainID: body.DomainID,
		Name:     body.DistributorName,
		UserID:   userId,
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
	// if err := valkeyCache.Revalidate("distributors:" + "1" + ":" + userId); err != nil {
	// 	log.Printf("Error in revalidating distributors cache: %s", err)
	// }
	return c.Status(201).JSON(fiber.Map{
		"msg": "distributor created successfully",
	})

}

func UpdateDistributor(c *fiber.Ctx) error {

	distributorId := c.Params("distributorId")
	userId := c.Params("userId")

	type exisitingDistributor struct {
		Name     string `json:"distributor_name"`
		DomainId string `json:"domain_id"`
	}

	body := new(exisitingDistributor)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}
	if result := database.DbConn.Model(&model.Distributor{}).Where("id=? AND user_id=?", distributorId, userId).Updates(
		model.Distributor{
			Name:     body.Name,
			DomainID: body.DomainId,
		},
	); result.Error != nil {
		log.Printf("Error updating distributor %s", result.Error.Error())
		return c.Status(500).JSON("Error updating distributor")
	}
	// if err := valkeyCache.Revalidate("distributors:" + userId); err != nil {
	// 	log.Printf("Error in revalidating distributors cache: %s", err)
	// }
	return c.Status(200).JSON("distributor updated successfully")
}

func DeleteDistributor(c *fiber.Ctx) error {
	distributorId := c.Params("distributorId")
	userId := c.Params("userId")
	var pgErr *pgconn.PgError
	if result := database.DbConn.Where("id=? AND user_id=?", distributorId, userId).Delete(&model.Distributor{}); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No distributor found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No Record found")
		}
		if errors.As(result.Error, &pgErr) {
			if pgErr.Code == "23503" {
				log.Printf("Error:Delete associated domains and bills to delete distributor. %s", result.Error.Error())
				return c.Status(400).JSON("Error:Delete associated domains and bills to delete distributor")
			}
		}
		log.Printf("Error deleting Distributor %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting Distributor")
	}
	// if err := valkeyCache.Revalidate("distributors:" + userId); err != nil {
	// 	log.Printf("Error in revalidating distributors cache: %s", err)
	// }
	return c.Status(200).JSON("Distributor deleted successfully")
}

func GetDistributor(c *fiber.Ctx) error {

	distributorId := c.Params("distributorId")
	userId := c.Params("userId")

	var distributor model.Distributor

	if err := database.DbConn.Limit(1).Where("id =? AND user_id =?", distributorId, userId).Find(&distributor).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error " + err.Error(),
		})
	}

	return c.Status(200).JSON(distributor)

}
