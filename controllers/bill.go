package controller

import (
	// "encoding/json"
	// "github.com/devMukulSingh/billManagementServer.git/db"
	// "time"
	// "github.com/jackc/pgx/v5/pgtype"
	// "github.com/devMukulSingh/billManagementServer.git/model"
	// "encoding/json"
	// "errors"
	"log"

	"github.com/devMukulSingh/billManagementServer.git/database"
	"github.com/devMukulSingh/billManagementServer.git/dbConnection"
	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	// "github.com/jackc/pgx/v5/pgtype"
	// "github.com/go-playground/validator/v10"
	// "github.com/devMukulSingh/billManagementServer.git/valkeyCache"
	// "gorm.io/gorm/clause"
)

func GetBills(c *fiber.Ctx) error {

	userId := c.Params("userId")
	var queryParams types.Query

	if err := c.QueryParser(&queryParams); err != nil {
		log.Print(err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in parsing query params " + err.Error(),
		})
	}

	data, err := dbconnection.Queries.GetBills(dbconnection.Ctx, database.GetBillsParams{
		UserID: userId,
		Offset: queryParams.Page,
		Limit:  queryParams.Limit,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in getting Bills " + err.Error(),
		})
	}
	// cache, err := valkeyCache.GetValue("bills:" + userId)
	// if err != nil {
	// 	if err.Error() != "valkey nil message" {
	// 		log.Printf("Error in getting cached bills : %s", err)
	// 	}
	// } else {
	// 	c.Set("Content-Type", "application/json")
	// 	return c.Status(200).SendString(cache)
	// }
	// if err := valkeyCache.SetValue("bills:"+userId, jsonBills); err != nil {
	// 	log.Printf("Error in setting value in valkey %s ", err.Error())
	// }
	return c.Status(200).JSON(data)

}

func PostBill(c *fiber.Ctx) error {

	body := new(types.Bill)
	userId := c.Params("userId")

	if err := c.BodyParser(body); err != nil {
		log.Printf("error parsing request body %s", err.Error())
		return c.Status(400).JSON("Error :error parsing request body")
	}

	//TODO: execute all queries inside transaction
	billId, err := dbconnection.Queries.PostBill(dbconnection.Ctx, database.PostBillParams{
		Date:          body.Date,
		TotalAmount:   body.TotalAmount,
		IsPaid:        body.IsPaid,
		UserID:        userId,
		DistributorID: body.DistributorId,
		DomainID:      body.DomainId,
	})
	if err != nil {
		log.Print(err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in posting bill " + err.Error(),
		})
	}
	n := len(body.BillItems)
	ids := make([]string,n);
	quantities := make([]int32,n);
	amounts := make([]int32,n);
	productIds := make([]string,n);
	billIds := make([]string,n);
	
	for i, billItem := range body.BillItems {
		quantities[i] = billItem.Quantity
		amounts[i] = billItem.Amount
		productIds[i] = billItem.ProductID
		billIds[i] = billId
		ids[i] = uuid.NewString()
	}
	
	if err := dbconnection.Queries.BatchInsertBillItems(dbconnection.Ctx,database.BatchInsertBillItemsParams{
		Ids: ids,
		Quantities: quantities,
		Amounts: amounts,
		Productids: productIds,
		Billids: billIds,
	} ); err != nil {
		log.Print(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in posting bill " + err.Error(),
		})
	}
	// billItems,err := json.Marshal(body.BillItems);
	// if err!=nil{
	// 		log.Print(err)
	// }
	// if err:= dbconnection.Queries.BatchInsertBillItems(dbconnection.Ctx,billItems);
	// err!= nil{
	// 	log.Print(err)
	// }

	// }
	// if err := valkeyCache.Revalidate("bills:" + userId); err != nil {
	// 	log.Printf("Error in revalidating bills cache: %s", err)
	// }
	return c.Status(201).JSON(fiber.Map{
		"msg": "bill created successfully",
	})

}

func UpdateBill(c *fiber.Ctx) error {

	billId := c.Params("billId")
	userId := c.Params("userId")

	body := new(types.Bill)
	if err := c.BodyParser(body); err != nil {
		log.Printf("error parsing request body %s", err.Error())
		return c.Status(400).JSON("Error :error parsing request body")
	}

	//TODO: execute all queries inside transaction
	if err := dbconnection.Queries.UpdateBill(dbconnection.Ctx, database.UpdateBillParams{
		Date:          body.Date,
		TotalAmount:   body.TotalAmount,
		IsPaid:        body.IsPaid,
		UserID:        userId,
		DistributorID: body.DistributorId,
		DomainID:      body.DomainId,
	});err != nil {
		log.Print(err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in posting bill " + err.Error(),
		})
	}

	n := len(body.BillItems)
	ids := make([]string,n);
	quantities := make([]int32,n);
	amounts := make([]int32,n);
	productIds := make([]string,n);
	billIds := make([]string,n);
	
	for i, billItem := range body.BillItems {
		quantities[i] = billItem.Quantity
		amounts[i] = billItem.Amount
		productIds[i] = billItem.ProductID
		billIds[i] = billId
		ids[i] = uuid.NewString()
	}

	if err:= dbconnection.Queries.DeleteManyBillItems(dbconnection.Ctx,billId);err!=nil{
		log.Print(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error":"Error deleting bill items : "+ err.Error(),
		})
	}
	
	if err := dbconnection.Queries.BatchInsertBillItems(dbconnection.Ctx,database.BatchInsertBillItemsParams{
		Ids: ids,
		Quantities: quantities,
		Amounts: amounts,
		Productids: productIds,
		Billids: billIds,
	} ); err != nil {
		log.Print(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in posting bill " + err.Error(),
		})
	}

	// if err := valkeyCache.Revalidate("bills:" + userId); err != nil {
	// 	log.Printf("Error in revalidating bills cache: %s", err)
	// }
	return c.Status(200).JSON("bill updated successfully")
}

func DeleteBill(c *fiber.Ctx) error {

	var params types.BillParams
	if err := c.ParamsParser(&params); err != nil {
		log.Print(err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in parsing params :" + err.Error(),
		})
	}
	if err := c.ParamsParser(&params); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in deleting bills " + err.Error(),
		})
	}

	if err := dbconnection.Queries.DeleteBill(dbconnection.Ctx, database.DeleteBillParams{
		ID:     params.BillId,
		UserID: params.UserId,
	}); err != nil {
		log.Print(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in deleting bill " + err.Error(),
		})
	}

	// if err := valkeyCache.Revalidate("bills:" + userId); err != nil {
	// 	log.Printf("Error in revalidating bills cache: %s", err)
	// }
	return c.Status(200).JSON("Bill deleted successfully")
}

// func GetBill(c *fiber.Ctx) error {

// 	userId := c.Params("billId")
// 	var params types.BillParams;
// 	if err:= c.ParamsParser(&params);err!= nil{
// 		log.Print(err)
// 		return c.Status(400).JSON(fiber.Map{
// 			"error":"Error in parsing params " + err.Error(),
// 		})
// 	}

// 	if err:= c.QueryParser(&types.Query{});err!=nil{
// 		return c.Status(400).JSON(fiber.Map{
// 				"error":"Error in parsing query params "+ err.Error(),
// 			})
// 	}

// 	data,err := dbconnection.Queries.GetBills(dbconnection.Ctx,userId);
// 	if err != nil{
// 		return c.Status(400).JSON(fiber.Map{
// 			"error":"Error in getting Bills " + err.Error(),
// 		})
// 	}
// 	// var bill model.Bill

// 	// if err := database.DbConn.Limit(1).Where("id =? AND user_id=?", billId, userId).Find(&bill).Error; err != nil {
// 	// 	return c.Status(500).JSON(fiber.Map{
// 	// 		"error": "Internal server error " + err.Error(),
// 	// 	})
// 	// }
// 	return c.Status(200).JSON(bill)
// }
