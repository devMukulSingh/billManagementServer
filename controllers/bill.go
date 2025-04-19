package controller

import (
	"fmt"
	"log"
	"github.com/devMukulSingh/billManagementServer.git/database"
	"github.com/devMukulSingh/billManagementServer.git/dbConnection"
	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	// "github.com/devMukulSingh/billManagementServer.git/valkeyCache"
)

func GetSearchedBills(c *fiber.Ctx) error {
	userId := c.Params("userId")
	var queries types.SearchBillsQuery
	if err := c.QueryParser(&queries); err != nil {
		log.Print(err)
	}

	data, err := dbconnection.Queries.GetSearchedBills(dbconnection.Ctx, database.GetSearchedBillsParams{
		UserID:      userId,
		Offset:      (queries.Page - 1) * queries.Limit,
		Limit:       queries.Limit,
		CreatedAt:   queries.StartDate,
		CreatedAt_2: queries.EndDate,
	})
	if err != nil {
		log.Print(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in getting searched bills : " + err.Error(),
		})
	}
	count, err := dbconnection.Queries.GetSearchedBillsCount(dbconnection.Ctx, database.GetSearchedBillsCountParams{
		UserID:      userId,
		CreatedAt:   queries.StartDate,
		CreatedAt_2: queries.EndDate,
		
	})
	if err != nil {
		log.Print(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in getting searched bills : " + err.Error(),
		})
	}

	type Response struct{
		Data		[]database.GetSearchedBillsRow		`json:"data"`
		Count 		int64			`json:"count"`
	}
	return c.Status(200).JSON(Response{
		Data: data,
		Count: count,
	})

}

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
		Offset: (queryParams.Page - 1) * queryParams.Limit,
		Limit:  queryParams.Limit,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in getting Bills " + err.Error(),
		})
	}
	count, err := dbconnection.Queries.GetBillsCount(dbconnection.Ctx, userId)
	if err != nil {
		log.Print(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in getting counts :" + err.Error(),
		})
	}
	type Response struct {
		Data  []database.GetBillsRow `json:"data"`
		Count int64                  `json:"count"`
	}
	response := Response{
		Data:  data,
		Count: count,
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
	return c.Status(200).JSON(response)

}

func PostBill(c *fiber.Ctx) error {

	body := new(types.Bill)
	userId := c.Params("userId")
	c.BodyParser(body)

	tx, err := dbconnection.Connection.Begin(dbconnection.Ctx)
	if err != nil {
		log.Print(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Unable to begin transaction : " + err.Error(),
		})
	}
	defer tx.Rollback(dbconnection.Ctx)

	qtx := database.New(tx)
	billId, err := qtx.PostBill(dbconnection.Ctx, database.PostBillParams{
		Date:          body.Date,
		TotalAmount:   body.TotalAmount,
		IsPaid:        body.IsPaid,
		UserID:        userId,
		DistributorID: body.DistributorId,
		DomainID:      body.DomainId,
		
	})
	if err != nil {
		return fmt.Errorf("failed to create bill: %w", err)
	}
	n := len(body.BillItems)
	ids := make([]string, n)
	quantities := make([]int32, n)
	amounts := make([]int32, n)
	productIds := make([]string, n)
	billIds := make([]string, n)

	for i, billItem := range body.BillItems {
		quantities[i] = billItem.Quantity
		amounts[i] = billItem.Amount
		productIds[i] = billItem.ProductID
		billIds[i] = billId
		ids[i] = uuid.NewString()
	}

	if err := qtx.BatchInsertBillItems(dbconnection.Ctx, database.BatchInsertBillItemsParams{
		Ids:        ids,
		Quantities: quantities,
		Amounts:    amounts,
		Productids: productIds,
		Billids:    billIds,
	}); err != nil {
		return fmt.Errorf("failed to create bill: %w", err)
	}

	if err := tx.Commit(dbconnection.Ctx); err != nil {
		log.Print(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in posting bill " + err.Error(),
		})
	}

	// if err := valkeyCache.Revalidate("bills:" + userId); err != nil {
	// 	log.Printf("Error in revalidating bills cache: %s", err)
	// }
	return c.Status(201).JSON(fiber.Map{
		"msg": "bill created successfully",
	})

}

func UpdateBill(c *fiber.Ctx) error {

	var params types.BillParams
	if err := c.ParamsParser(&params); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Error parsing params :" + err.Error(),
		})
	}
	body := new(types.Bill)
	c.BodyParser(body)

	tx, err := dbconnection.Connection.Begin(dbconnection.Ctx)
	if err != nil {
		log.Print(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Unable to begin transaction : " + err.Error(),
		})
	}
	defer tx.Rollback(dbconnection.Ctx)
	qtx := database.New(tx)

	if err := qtx.UpdateBill(dbconnection.Ctx, database.UpdateBillParams{
		ID:            params.BillId,
		Date:          body.Date,
		TotalAmount:   body.TotalAmount,
		IsPaid:        body.IsPaid,
		UserID:        params.UserId,
		DistributorID: body.DistributorId,
		DomainID:      body.DomainId,
	}); err != nil {
		return fmt.Errorf("error in updating bill:  %w", err)
	}

	n := len(body.BillItems)
	ids := make([]string, n)
	quantities := make([]int32, n)
	amounts := make([]int32, n)
	productIds := make([]string, n)
	billIds := make([]string, n)

	for i, billItem := range body.BillItems {
		quantities[i] = billItem.Quantity
		amounts[i] = billItem.Amount
		productIds[i] = billItem.ProductID
		billIds[i] = params.BillId
		ids[i] = uuid.NewString()
	}

	if err := qtx.DeleteManyBillItems(dbconnection.Ctx, params.BillId); err != nil {
		return fmt.Errorf("error deleting bill items: %w", err)
	}

	if err := qtx.BatchInsertBillItems(dbconnection.Ctx, database.BatchInsertBillItemsParams{
		Ids:        ids,
		Quantities: quantities,
		Amounts:    amounts,
		Productids: productIds,
		Billids:    billIds,
	}); err != nil {
		log.Print(err.Error())
		return fmt.Errorf("error in updating bill %w", err)
	}
	if err := tx.Commit(dbconnection.Ctx); err != nil {
		log.Print(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in updating bill " + err.Error(),
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
