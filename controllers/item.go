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

func GetAllItems(c *fiber.Ctx) error {

	userId := c.Params("userId")

	cachedItems,err := valkeyCache.GetValue("billItems:"+userId);
	if err!=nil{
		if err.Error()!="valkey nil message"{
			log.Printf("Error in getting cached item from valkey: %s",err);
		}
	}else{
		c.Set("Content-Type","application/json")
		return c.SendString(cachedItems)
	}

	var items []model.Item

	if err := database.DbConn.Where("user_id =?", userId).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error " + err.Error(),
		})
	}
	jsonItems ,err := json.Marshal(items);

	if err!=nil{
		log.Printf("error converting to json, items : %s",err)
	} 
	if err := valkeyCache.SetValue("billItems:"+userId,jsonItems);err!=nil{
		log.Printf("error setting billItems in valkey : %s",err);
	}
	return c.Status(200).JSON(items)
}

func PostItem(c *fiber.Ctx) error {
	body := new(types.Item)
	userId := c.Params("userId")
	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}

	result := database.DbConn.Create(&model.Item{
		Name:   body.Name,
		Rate:   body.Rate,
		UserID: userId,
	})
	if result.Error != nil {
		log.Printf("Error in saving Items into db %s", result.Error.Error())
		return c.Status(500).JSON("Internal server error")
	}
	if err := valkeyCache.Revalidate("billItems:"+userId);err!=nil{
		log.Printf("Error in revalidating 'items' from valkey: %s",err)
	}
	return c.Status(201).JSON(fiber.Map{
		"msg": "Item created successfully",
	})
}

func UpdateItem(c *fiber.Ctx) error {
	itemId := c.Params("itemId")
	userId := c.Params("user_id")

	body := new(types.Item)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}
	if result := database.DbConn.Where("id=? AND item_id=?",itemId,userId).Model(&model.Item{}).Updates(
		model.Item{
			Name:   body.Name,
			Rate:   body.Rate,
			UserID: body.UserID,
		},
	); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No item found %s", result.Error.Error())
			return c.Status(400).JSON("No item found")
		}
		log.Printf("Error updating item %s", result.Error.Error())
		return c.Status(500).JSON("Error updating item")
	}
	if err := valkeyCache.Revalidate("billItems:"+userId);err!=nil{
		log.Printf("Error in revalidating 'billItems' from valkey: %s",err)
	}
	
	return c.Status(200).JSON("item updated successfully")
}

func DeleteItem(c *fiber.Ctx) error {
	itemId := c.Params("itemId")
	userId := c.Params("userId")

	if result := database.DbConn.Where("id=? AND user_id=?", itemId, userId).
		Delete(&model.Item{}); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No item found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No Record found")
		}
		log.Printf("Error deleting item %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting item")
	}
	if err := valkeyCache.Revalidate("billItems:"+userId);err!=nil{
		log.Printf("Error in revalidating 'billItems' from valkey: %s",err)
	}
	return c.Status(200).JSON("item deleted successfully")
}
