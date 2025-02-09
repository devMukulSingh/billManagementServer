package controller

import (
	"log"
	// "time"
	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
)



func PostBillController(c *fiber.Ctx) error {
	type Item struct {
		Name     				string `json:"name"`
		Rate     				int    `json:"rate"`
		Amount   				int    `json:"amount"`
		Quantity 				int    `json:"quantity"`
	}
	type Bill struct {
		DistributorId 			string `json:"distributor_id"`
		DomainId      			string `json:"domain_id"`
		Date        			string `json:"date"`
		IsPaid      			bool   `json:"isPaid"`
		Items       			[]Item `json:"items"`
		TotalAmount 			int    `json:"totalAmount"`
	}

	body := new(Bill)
	if err := c.BodyParser(body); err != nil {
		log.Printf("error parsing request body %s", err.Error())
		return c.Status(400).JSON("Error :error parsing request body")
	}
	log.Print(body)
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
		Distributor: model.Distributor{
			DomainID: body.DomainId,
			Base: model.Base{
				ID: body.DistributorId,
			},
		},
		Items: items,
		IsPaid: body.IsPaid,
		Date: body.Date,
		TotalAmount: body.TotalAmount,
		Domain: model.Domain{
			Base: model.Base{
				ID: body.DomainId,
			},
		},
	})
	if result.Error != nil {
		log.Printf("Error saving into db %s", result.Error.Error())
		return c.Status(500).JSON("Internal server error")
	}
	return c.Status(201).JSON("bill created successfully")
}

