package controller

import (
	"log"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UpdateBill(c *fiber.Ctx) error {

	id := c.Params("id")

	var existingBill model.Bill
	if result := database.DbConn.First(&existingBill, "id = ?", id); result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			log.Print("Failed to find bill with particular id")
			return c.Status(400).JSON(fiber.Map{
				"error": "Bill not found", 
			}) 
		}
		return c.Status(500).JSON("Internal server error")
	}
	body := new(Bill)
	if err := c.BodyParser(body); err != nil {
		log.Printf("error parsing request body %s", err.Error())
		return c.Status(400).JSON("Error :error parsing request body")
	}

	if err := database.DbConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).
			Model(&model.Bill{}).Where("id=?", id).
			Updates(body).Error; err != nil {
			return err
		}
		return nil
	}); err!= nil{
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
