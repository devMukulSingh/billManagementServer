package controller

import (
	"errors"
	"log"
	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/model"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func DeleteDistributor(c *fiber.Ctx) error {
	id := c.Params("id")
	var pgErr *pgconn.PgError
	if result := database.DbConn.Delete(&model.Distributor{}, "id =?", id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No Distributor found %s", result.Error.Error())
			return c.Status(400).JSON("Error:No Distributor found")
		}
		if errors.As(result.Error, &pgErr) {
		if pgErr.Code=="23503"  {
			log.Printf("Error:Delete associated domains and bills to delete distributor. %s",result.Error.Error())
			return c.Status(400).JSON("Error:Delete associated domains and bills to delete distributor")
		}
		}
		log.Printf("Error deleting Distributor %s", result.Error.Error())
		return c.Status(500).JSON("Error deleting Distributor")
	}

	return c.Status(200).JSON("Distributor deleted successfully")
}
