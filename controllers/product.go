package controller

import (
	// "encoding/json"
	"errors"
	"log"
	"time"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/devMukulSingh/billManagementServer.git/valkeyCache"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllProducts(c *fiber.Ctx) error {

	userId := c.Params("userId")

	// cachedItems,err := valkeyCache.GetValue("billItems:"+userId);
	// if err!=nil{
	// 	if err.Error()!="valkey nil message"{
	// 		log.Printf("Error in getting cached item from valkey: %s",err);
	// 	}
	// }else{
	// 	c.Set("Content-Type","application/json")
	// 	return c.SendString(cachedItems)
	// }
	var count int64;
	type Product struct {
		ID        		string          `json:"id" `
		CreatedAt 		time.Time       `json:"created_at"`
		Name      		string          `json:"name" `
		Rate    		int 			`json:"rate"`
	}
	var products []Product

	if err := database.DbConn.
	Model(&model.Product{}).
	Count(&count).
	Select("name","rate","id","created_at",).
	Where("user_id =?", userId).
	Scan(&products).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error " + err.Error(),
		})
	}

	// jsonItems ,err := json.Marshal(items);

	// if err!=nil{
	// 	log.Printf("error converting to json, items : %s",err)
	// } 
	// if err := valkeyCache.SetValue("billItems:"+userId,jsonItems);err!=nil{
	// 	log.Printf("error setting billItems in valkey : %s",err);
	// }
	type Response struct{
		Data		[]Product			`json:"data"`
		Count		int64				`json:"count"`
	}
	response := Response{
		Data: products,
		Count: count,
	}
	return c.Status(200).JSON(response)
}
func GetProducts(c *fiber.Ctx) error {

	userId := c.Params("userId")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	// cachedItems,err := valkeyCache.GetValue("billItems:"+userId);
	// if err!=nil{
	// 	if err.Error()!="valkey nil message"{
	// 		log.Printf("Error in getting cached item from valkey: %s",err);
	// 	}
	// }else{
	// 	c.Set("Content-Type","application/json")
	// 	return c.SendString(cachedItems)
	// }
	var count int64;
	type Product struct {
		ID        		string          `json:"id" `
		CreatedAt 		time.Time       `json:"created_at"`
		Name      		string          `json:"name" `
		Rate    		int 			`json:"rate"`
	}
	var products []Product

	if err := database.DbConn.
	Model(&model.Product{}).
	Count(&count).
	Select("name","rate","id","created_at",).
	Limit(limit).
	Offset((page-1) * limit).
	Where("user_id =?", userId).
	Scan(&products).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error " + err.Error(),
		})
	}

	// jsonItems ,err := json.Marshal(items);

	// if err!=nil{
	// 	log.Printf("error converting to json, items : %s",err)
	// } 
	// if err := valkeyCache.SetValue("billItems:"+userId,jsonItems);err!=nil{
	// 	log.Printf("error setting billItems in valkey : %s",err);
	// }
	type Response struct{
		Data		[]Product			`json:"data"`
		Count		int64			`json:"count"`
	}
	response := Response{
		Data: products,
		Count: count,
	}
	return c.Status(200).JSON(response)
}
func PostProduct(c *fiber.Ctx) error {
	body := new(types.Product)
	userId := c.Params("userId")
	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}

	result := database.DbConn.Create(&model.Product{
		Name:   body.Name,
		Rate:   body.Rate,
		UserID: userId,
	})
	if result.Error != nil {
		log.Printf("Error in saving Products into db %s", result.Error.Error())
		return c.Status(500).JSON("Internal server error")
	}
	if err := valkeyCache.Revalidate("billItems:"+userId);err!=nil{
		log.Printf("Error in revalidating 'Products' from valkey: %s",err)
	}
	return c.Status(201).JSON(fiber.Map{
		"msg": "Products created successfully",
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	productId := c.Params("productId")
	userId := c.Params("user_id")

	body := new(types.Product)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}
	if result := database.DbConn.Where("id=? AND product_id=?",productId,userId).Model(&model.Product{}).Updates(
		model.Product{
			Name:   body.Name,
			Rate:   body.Rate,
			UserID: body.UserID,
		},
	); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No Product found %s", result.Error.Error())
			return c.Status(400).JSON("No Product found")
		}
		log.Printf("Error updating product %s", result.Error.Error())
		return c.Status(500).JSON("Error updating Product")
	}
	// if err := valkeyCache.Revalidate("billItems:"+userId);err!=nil{
	// 	log.Printf("Error in revalidating 'billItems' from valkey: %s",err)
	// }
	
	return c.Status(200).JSON("Product updated successfully")
}

func DeleteProduct(c *fiber.Ctx) error {
	productId := c.Params("productId")
	userId := c.Params("userId")

	if result := database.DbConn.Where("id=? AND user_id=?", productId, userId).
		Delete(&model.Product{}); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No Product found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No Record found")
		}
		log.Printf("Error deleting Product %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting Product")
	}
	// if err := valkeyCache.Revalidate("billItems:"+userId);err!=nil{
	// 	log.Printf("Error in revalidating 'billItems' from valkey: %s",err)
	// }
	return c.Status(200).JSON("Product deleted successfully")
}
