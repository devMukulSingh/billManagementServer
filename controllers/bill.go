package controller

import (
	"errors"
	"log"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAllBills(c *fiber.Ctx) error {

	userId := c.Params("userId")

	var bills []model.Bill

	if err := database.DbConn.Where("user_id =?", userId).Find(&bills).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error " + err.Error(),
		})
	}

	return c.Status(200).JSON(bills)

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

	var items []model.Item
	for _, itemReq := range body.Items {
		item := model.Item{
			Name:     itemReq.Name,
			Rate:     itemReq.Rate,
			Amount:   itemReq.Amount,
			Quantity: itemReq.Quantity,
		}
		items = append(items, item)
	}

	result := database.DbConn.Create(&model.Bill{
		UserID: userId, 
		DistributorID: body.DistributorId,
		DomainID: body.DomainId,
		Items:       items,
		IsPaid:      body.IsPaid,
		Date:        body.Date,
		TotalAmount: body.TotalAmount,
	})
	if result.Error != nil {

		log.Printf("Error saving into db %s", result.Error.Error())

		return c.Status(500).JSON("Internal server error")
	}
	return c.Status(201).JSON(fiber.Map{
		"msg": "bill created successfully",
	})
}

func UpdateBill(c *fiber.Ctx) error {

	billId := c.Params("billId")

	var existingBill model.Bill
	if result := database.DbConn.First(&existingBill, "id = ?", billId); result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			log.Print("Failed to find bill with particular id")
			return c.Status(400).JSON(fiber.Map{
				"error": "Bill not found",
			})
		}
		return c.Status(500).JSON("Internal server error")
	}
	body := new(types.Bill)
	if err := c.BodyParser(body); err != nil {
		log.Printf("error parsing request body %s", err.Error())
		return c.Status(400).JSON("Error :error parsing request body")
	}

	if err := database.DbConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).
			Model(&model.Bill{}).Where("id=?", billId).
			Updates(body).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Print("Internal server errror")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update bill",
		})
	}

	// var items []model.Item
	// for _, itemReq := range body.Items {
	// 	item := model.Item{
	// 		Base: model.Base{
	// 			ID: itemReq.ID,
	// 		},
	// 		Name:     itemReq.Name,
	// 		Rate:     itemReq.Rate,
	// 		Amount:   itemReq.Amount,
	// 		Quantity: itemReq.Quantity,
	// 	}

	// 	items = append(items, item)
	// }

	// result := database.DbConn.Update

	return c.Status(200).JSON("bill updated successfully")
}

func DeleteBill(c *fiber.Ctx) error {
	billId := c.Params("billId")
	var existingBill model.Bill

	if result := database.DbConn.First(&existingBill, "id=?", billId); result != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No bill found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No Record found")
		}
	}

	if result := database.DbConn.Select(clause.Associations).Delete(&existingBill); result.Error != nil {
		log.Printf("Error deleting Bill %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting Bill")
	}

	return c.Status(200).JSON("Bill deleted successfully")
}
