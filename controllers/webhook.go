package controller

import (
	"encoding/json"
	"log"
	"os"

	database "github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
	svix "github.com/svix/svix-webhooks/go"
)

type Event struct {
	Data Data   `json:"data"`
	Type string `json:"type"`
}
type Data struct {
	First_name      string          `json:"first_name"`
	Last_name       string          `json:"last_name"`
	Id              string          `json:"id"`
	Email_Addresses []Email_Address `json:"email_addresses"`
}
type Email_Address struct {
	Email_Address string `json:"email_address"`
}

func Webhook(c *fiber.Ctx) error {
	secret := os.Getenv("SIGNING_SECRET")
	body := c.Body()
	headers := c.GetReqHeaders()

	var event Event
	if err := json.Unmarshal(body, &event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON payload",
		})
	}
	if event.Type == "user.created" {

		wh, err := svix.NewWebhook(secret)
		if err != nil {
			log.Printf(`Error in creating webhook %s`, err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON("Invalid signature")
		}

		err = wh.Verify(body, headers)
		if err != nil {
			log.Printf(`Error in verifying headers %s`, err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON("Invalid signature")
		}

		if result := database.DbConn.Create(&model.User{
			Name:  event.Data.First_name + " " + event.Data.Last_name,
			Email: event.Data.Email_Addresses[0].Email_Address,
			ID: event.Data.Id,
		}); result.Error != nil {
			log.Printf("failed to create user into db %s", result.Error.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to create user into db " + result.Error.Error(),
			})
		}
		return c.Status(201).JSON("User creatd successfully")
	}
	return c.Status(200).JSON("Other event than user.created")
}
