package controller

import (
	"encoding/json"
	"errors"
	"log"
	"time"
	// "github.com/go-playground/validator/v10"
	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/devMukulSingh/billManagementServer.git/valkeyCache"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
)

func GetAllBills(c *fiber.Ctx) error {

	userId := c.Params("userId")
	cache, err := valkeyCache.GetValue("bills:" + userId)
	if err != nil {
		if err.Error() != "valkey nil message" {
			log.Printf("Error in getting cached bills : %s", err)
		}
	} else {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString(cache)
	}

	type Data struct {
		ID          string          `json:"id" gorm:"type:uuid;primary_key;"`
		CreatedAt   time.Time       `json:"created_at"`
		Date        time.Time       `json:"date"`
		IsPaid      bool            `json:"is_paid"`
		TotalAmount int             `json:"total_amount"`
		Distributor json.RawMessage `json:"distributor"`
		Domain      json.RawMessage `json:"domain"`
		BillItems   json.RawMessage `json:"bill_items"`
	}

	var data []Data
	if err := database.DbConn.Model(&model.Bill{}).
		Joins("JOIN distributors ON distributors.id = bills.distributor_id").
		Joins("JOIN domains ON domains.id = bills.domain_id").
		Select(`
		bills.id,
		bills.created_at,
	    bills.date,
	    bills.is_paid,
	    bills.total_amount,
		row_to_json(distributors) as distributor,
		row_to_json(domains) as domain,
	    (
	      SELECT json_agg(
	        json_build_object(
	          'id', bi.id,
	          'quantity', bi.quantity,
	          'amount', bi.amount,
	          'item', json_build_object(
	            'id', it.id,
	            'name', it.name,
	            'rate', it.rate
	          )
	        )
	      )
	      FROM bill_items AS bi
	      JOIN items AS it ON it.id = bi.item_id
	      WHERE bi.bill_id = bills.id
	    ) as bill_items
	`).
		Where("bills.user_id =?", userId).
		Scan(&data).Error; err != nil {

		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error " + err.Error(),
		})
	}
	jsonBills, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error in converting to json: %s", err)
	}
	if err := valkeyCache.SetValue("bills:"+userId, jsonBills); err != nil {
		log.Printf("Error in setting value in valkey %s ", err.Error())
	}
	return c.Status(200).JSON(data)

}

func GetBill(c *fiber.Ctx) error {

	userId := c.Params("billId")
	billId := c.Params("userId")

	var bill model.Bill

	if err := database.DbConn.Limit(1).Where("id =? AND user_id=?", billId, userId).Find(&bill).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error " + err.Error(),
		})
	}
	return c.Status(200).JSON(bill)
}

func PostBill(c *fiber.Ctx) error {

	body := new(types.Bill)
	userId := c.Params("userId")
	if err := c.BodyParser(body); err != nil {
		log.Printf("error parsing request body %s", err.Error())
		return c.Status(400).JSON("Error :error parsing request body")
	}

	bills := model.Bill{
		UserID:        userId,
		DistributorID: body.DistributorId,
		DomainID:      body.DomainId,
		IsPaid:        body.IsPaid,
		Date:          body.Date,
		TotalAmount:   body.TotalAmount,
		// Items
	}
	if err := database.DbConn.Create(&bills).Error; err != nil {
		log.Printf("Error saving into db %s", err.Error())
		return c.Status(500).JSON("Internal server error")
	}

	var items []model.BillItem
	for _, itemReq := range body.BillItems {
		item := model.BillItem{
			BillID:   bills.Base.ID,
			ItemID:   itemReq.ItemID,
			Amount:   itemReq.Amount,
			Quantity: itemReq.Quantity,
		}
		items = append(items, item)
	}

	if err := database.DbConn.Create(&items).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Printf("Record already exists,try another - %s", err.Error())
			return c.Status(409).JSON(fiber.Map{
				"error": "Item already exists, try another",
			})
		}
	}
	if err := valkeyCache.Revalidate("bills:" + userId); err != nil {
		log.Printf("Error in revalidating bills cache: %s", err)
	}
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
	//todo: only save altered values
	if err := database.DbConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).
			Model(&model.Bill{}).Where("id=? AND user_id", billId,userId).
			Updates(body).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Print("Failed to find bill with particular id")
			return c.Status(400).JSON(fiber.Map{
				"error": "Bill not found",
			})
		}
		log.Print("Internal server errror")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update bill",
		})
	}
	if err := valkeyCache.Revalidate("bills:" + userId); err != nil {
		log.Printf("Error in revalidating bills cache: %s", err)
	}
	return c.Status(200).JSON("bill updated successfully")
}

func DeleteBill(c *fiber.Ctx) error {

	billId := c.Params("billId")
	userId := c.Params("userId")

	if result := database.DbConn.Where("id=?", billId).Delete(&model.Bill{}); result.Error != nil {
		log.Printf("Error deleting Bill %s", result.Error.Error())
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No bill found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No Record found")
		}
		return c.Status(500).JSON("Error deleting Bill")
	}
	if err := valkeyCache.Revalidate("bills:" + userId); err != nil {
		log.Printf("Error in revalidating bills cache: %s", err)
	}
	return c.Status(200).JSON("Bill deleted successfully")
}
