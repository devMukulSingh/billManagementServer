package controller

import (
	"errors"
	"log"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)


func GetAllDistributors(c * fiber.Ctx) error{

	userId := c.Params("userId")

	var distributors []model.Distributor;

	if err := database.DbConn.Where("user_id =?",userId).Find(&distributors).Error; err!=nil{
		return c.Status(500).JSON(fiber.Map{
			"error":"Internal server error " + err.Error(),
		})
	}

	return c.Status(200).JSON(distributors);

}

func GetDistributor(c * fiber.Ctx) error{

	distributorId := c.Params("distributorId")
	userId := c.Params("userId")

	var distributor model.Distributor;

	if err := database.DbConn.Limit(1).Where("id =? AND user_id =?", distributorId,userId).Find(&distributor).Error; err!=nil{
		return c.Status(500).JSON(fiber.Map{
			"error":"Internal server error " + err.Error(),
		}) 
	}

	return c.Status(200).JSON(distributor);

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
		Name: body.DistributorName,
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

func UpdateDistributor(c *fiber.Ctx) error {

	distributorId := c.Params("distributorId")
	userId := c.Params("userId")

	type exisitingDistributor struct{
		Name					string			`json:"name"`
		DomainId 				string			`json:"domain_id"`
	}

	body := new(exisitingDistributor)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}
	if result := database.DbConn.Model(&model.Distributor{}).Where("id=? AND user_id",distributorId,userId).Updates(
		model.Distributor{
			Name: body.Name,
			DomainID: body.DomainId,
		},
	); result.Error != nil {
		log.Printf("Error updating distributor %s", result.Error.Error())
		return c.Status(500).JSON("Error updating distributor")
	}

	return c.Status(200).JSON("distributor updated successfully")
}

func DeleteDistributor(c *fiber.Ctx) error {
	distributorId := c.Params("distributorId")
	var existingDistributor model.Distributor;

	if result := database.DbConn.First(&existingDistributor, "id=?", distributorId); result != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No distributor found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No Record found")
		}
	}

	var pgErr *pgconn.PgError
	if result := database.DbConn.Delete(&existingDistributor); result.Error != nil {
		if errors.As(result.Error, &pgErr) {
		if pgErr.Code=="23503"  {
			log.Printf("Error:Delete associated domains and bills to delete distributor. %s",result.Error.Error())
			return c.Status(400).JSON("Error:Delete associated domains and bills to delete distributor")
		}
		}
		log.Printf("Error deleting Distributor %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting Distributor")
	}

	return c.Status(200).JSON("Distributor deleted successfully")
}
