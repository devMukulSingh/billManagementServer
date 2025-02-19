package controller

import (
	"errors"
	"log"
	"gorm.io/gorm"
	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/gofiber/fiber/v2"
)

func GetAllItems(c * fiber.Ctx) error{
	
	billId := c.Params("billId")

	var items []model.Item

	if err := database.DbConn.Where("bill_id =?",billId).Find(&items).Error; err!=nil{
		return c.Status(500).JSON(fiber.Map{
			"error":"Internal server error " + err.Error(),
		})
	}

	return c.Status(200).JSON(items);
}

// func GetItem(c * fiber.Ctx) error{
	
// 	userId := c.Params("userId")
// 	itemId := c.Params("itemId")
	
// 	var items []model.Item

// 	if err := database.DbConn.Limit(1).Where("user_id =? AND id",userId,itemId).Find(&items).Error; err!=nil{
// 		return c.Status(500).JSON(fiber.Map{
// 			"error":"Internal server error " + err.Error(),
// 		})
// 	}

// 	return c.Status(200).JSON(items);
// }

func PostItem(c *fiber.Ctx) error {
	body := new(types.Item)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}

	result := database.DbConn.Create(&model.Item{
		Name: body.Name,
		Rate: body.Rate,
		Amount: body.Amount,
		Quantity: body.Quantity,
	})
	if result.Error != nil {
		log.Printf("Error in saving Items into db %s", result.Error.Error())
		return c.Status(500).JSON("Internal server error")
	}

	return c.Status(201).JSON(fiber.Map{
		"msg": "Item created successfully",
	})
}

func UpdateItem(c *fiber.Ctx) error {
	itemId := c.Params("itemId")

	var existingItem model.Item;

	if result := database.DbConn.First(&existingItem, "id = ?", itemId); result.Error != nil {
		log.Printf("No item found %s", result.Error.Error())
		return c.Status(400).JSON("No item found")
	}

	body := new(types.Item)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}
	if result := database.DbConn.Model(&existingItem).Updates(
		model.Item{
			Name:     body.Name,
			Rate:     body.Rate,
			Amount:   body.Amount,
			Quantity: body.Rate,
		},
	); result.Error != nil {
		log.Printf("Error updating item %s", result.Error.Error())
		return c.Status(500).JSON("Error updating item")
	}

	return c.Status(200).JSON("item updated successfully")
}

func DeleteItem(c *fiber.Ctx) error {
	itemId := c.Params("itemId")

	var existingItem model.Item

	if result := database.DbConn.First(&existingItem, "id=?", itemId); result != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No item found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No Record found")
		}
	}

	if result := database.DbConn.Delete(&existingItem); result.Error != nil {
		log.Printf("Error deleting item %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting item")
	}

	return c.Status(200).JSON("item deleted successfully")
}
