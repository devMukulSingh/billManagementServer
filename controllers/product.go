package controller

import (
	// "encoding/json"
	"log"
	// "strconv"
	// "errors"
	// "github.com/devMukulSingh/billManagementServer.git/model"
	// "github.com/devMukulSingh/billManagementServer.git/db"

	"github.com/devMukulSingh/billManagementServer.git/database"
	dbconnection "github.com/devMukulSingh/billManagementServer.git/dbConnection"
	"github.com/devMukulSingh/billManagementServer.git/types"
	// "github.com/devMukulSingh/billManagementServer.git/valkeyCache"
	"github.com/gofiber/fiber/v2"
	// "gorm.io/gorm"
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

	// var count int64;
	// var products []Product

	data, err := dbconnection.Queries.GetAllProducts(dbconnection.Ctx, userId)
	if err != nil {
		log.Print(err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in getting Products " + err.Error(),
		})
	}
	// count, err := dbconnection.Queries.GetProductsCount(dbconnection.Ctx, userId)
	// if err != nil {
	// 	log.Print(err.Error())
	// 	return c.Status(400).JSON(fiber.Map{
	// 		"error": "Error in getting Products " + err.Error(),
	// 	})
	// }
	// if err := database.DbConn.
	// Model(&model.Product{}).
	// Count(&count).
	// Select("name","rate","id","created_at",).
	// Where("user_id =?", userId).
	// Scan(&products).Error; err != nil {
	// 	return c.Status(500).JSON(fiber.Map{
	// 		"error": "Internal server error " + err.Error(),
	// 	})
	// }

	// jsonItems ,err := json.Marshal(items);

	// if err!=nil{
	// 	log.Printf("error converting to json, items : %s",err)
	// }
	// if err := valkeyCache.SetValue("billItems:"+userId,jsonItems);err!=nil{
	// 	log.Printf("error setting billItems in valkey : %s",err);
	// }
	// type Response struct{
	// 	Data		[]database.Product			`json:"data"`
	// 	Count		int64						`json:"count"`
	// }
	// response := Response{
	// 	Data: data,
	// 	Count: count,
	// }
	return c.Status(200).JSON(data)
}
func GetProducts(c *fiber.Ctx) error {

	userId := c.Params("userId")
	var query types.Query

	if err := c.QueryParser(&query); err != nil {
		log.Print(err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in parsing query " + err.Error(),
		})
	}

	data, err := dbconnection.Queries.GetProducts(dbconnection.Ctx, database.GetProductsParams{
		UserID: userId,
		Offset: (query.Page - 1) * query.Limit,
		Limit:  query.Limit,
	})
	if err != nil {
		log.Print(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in getting Prodocts " + err.Error(),
		})
	}

	count, err := dbconnection.Queries.GetProductsCount(dbconnection.Ctx, userId)
	if err != nil {
		log.Print(err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in getting Products " + err.Error(),
		})
	}
	// cachedItems,err := valkeyCache.GetValue("billItems:"+userId);
	// if err!=nil{
	// 	if err.Error()!="valkey nil message"{
	// 		log.Printf("Error in getting cached item from valkey: %s",err);
	// 	}
	// }else{
	// 	c.Set("Content-Type","application/json")
	// 	return c.SendString(cachedItems)

	// jsonItems ,err := json.Marshal(items);

	// if err!=nil{
	// 	log.Printf("error converting to json, items : %s",err)
	// }
	// if err := valkeyCache.SetValue("billItems:"+userId,jsonItems);err!=nil{
	// 	log.Printf("error setting billItems in valkey : %s",err);
	// }
	type Response struct{
		Data		[]database.Product			`json:"data"`
		Count		int64			`json:"count"`
	}
	response := Response{
		Data: data,
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

	if err := dbconnection.Queries.PostProduct(dbconnection.Ctx, database.PostProductParams{
		Name:   body.Name,
		Rate:   int32(body.Rate),
		UserID: userId,
	}); err != nil {
		log.Print(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in post Product " + err.Error(),
		})
	}

	// result := database.DbConn.Create(&model.Product{
	// 	Name:   body.Name,
	// 	Rate:   body.Rate,
	// 	UserID: userId,
	// })
	// if result.Error != nil {
	// 	log.Printf("Error in saving Products into db %s", result.Error.Error())
	// 	return c.Status(500).JSON("Internal server error")
	// }

	// if err := valkeyCache.Revalidate("billItems:"+userId);err!=nil{
	// 	log.Printf("Error in revalidating 'Products' from valkey: %s",err)
	// }

	return c.Status(201).JSON(fiber.Map{
		"msg": "Products created successfully",
	})
}

func UpdateProduct(c *fiber.Ctx) error {

	var params types.ProductParams
	if err := c.ParamsParser(&params); err != nil {
		log.Print(err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in parsing params " + err.Error(),
		})
	}

	body := new(types.Product)
	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}

	if err := dbconnection.Queries.UpdateProduct(dbconnection.Ctx, database.UpdateProductParams{
		ID:     params.ProductId,
		UserID: params.UserId,
		Name:   body.Name,
		Rate:   body.Rate,
	}); err != nil {
		log.Print(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in updating product " + err.Error(),
		})
	}
	
	// if err := valkeyCache.Revalidate("billItems:"+userId);err!=nil{
	// 	log.Printf("Error in revalidating 'billItems' from valkey: %s",err)
	// }

	return c.Status(200).JSON("Product updated successfully")
}

func DeleteProduct(c *fiber.Ctx) error {

	var params types.ProductParams
	if err := c.ParamsParser(&params); err != nil {
		log.Print(err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in parsing params " + err.Error(),
		})
	}
	if err := dbconnection.Queries.DeleteProduct(dbconnection.Ctx, database.DeleteProductParams{
		ID:     params.ProductId,
		UserID: params.UserId,
	}); err != nil {
		log.Print(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in deleting product " + err.Error(),
		})
	}


	// if err := valkeyCache.Revalidate("billItems:"+userId);err!=nil{
	// 	log.Printf("Error in revalidating 'billItems' from valkey: %s",err)
	// }
	return c.Status(200).JSON("Product deleted successfully")
}
