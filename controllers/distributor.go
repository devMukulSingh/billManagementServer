package controller

import (
	"errors"
	"log"
	"strings"

	"github.com/devMukulSingh/billManagementServer.git/database"
	"github.com/devMukulSingh/billManagementServer.git/dbConnection"
	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/jackc/pgx/v5/pgconn"

	// "github.com/devMukulSingh/billManagementServer.git/valkeyCache"
	"github.com/gofiber/fiber/v2"
)

func GetSearchedDistributors(c *fiber.Ctx) error {
	userId := c.Params("userId")
	type Query struct {
		Page  int32  `query:"page"`
		Limit int32  `query:"limit"`
		Name  string `query:"name"`
	}
	var query Query
	if err := c.QueryParser(&query); err != nil {
		log.Print(err)
	}
	data, err := dbconnection.Queries.GetSearchedDistributors(dbconnection.Ctx, database.GetSearchedDistributorsParams{
		Name:  "%" + strings.ToLower(query.Name) + "%",
		UserID: userId,
		Offset: (query.Page - 1)*query.Limit,
		Limit:  query.Limit,
	})
	if err != nil {
		log.Print(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in getting searched distributors :" + err.Error(),
		})
	}

	count, err := dbconnection.Queries.GetSearchedDistributorsCount(dbconnection.Ctx, database.GetSearchedDistributorsCountParams{
		Name:   "%" + strings.ToLower(query.Name) + "%",
		UserID: userId,
	})
	if err != nil {
		log.Print(err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in getting searched distributors count :" + err.Error(),
		})
	}

	type Response struct{
		Data		[]database.GetSearchedDistributorsRow		`json:"data"`
		Count 		int64			`json:"count"`
	}

	return c.Status(200).JSON(Response{
		Data: data,
		Count: count,
	})
}

func GetDistributors(c *fiber.Ctx) error {

	userId := c.Params("userId")
	page := int32(c.QueryInt("page", 1))
	limit := int32(c.QueryInt("limit", 10))
	// cached, err := valkeyCache.GetValue("distributors:" + page + ":" + userId)

	// if err != nil {
	// 	if err.Error() != "valkey nil message" {
	// 		log.Printf("Error in getting cached bills : %s", err)
	// 	}
	// } else {
	// 	c.Set("Content-Type", "application/json")
	// 	return c.SendString(cached)
	// }

	data, err := dbconnection.Queries.GetDistributors(dbconnection.Ctx, database.GetDistributorsParams{
		UserID: userId,
		Offset: (page - 1) * limit,
		Limit:  limit,
	})
	if err != nil {
		log.Print(err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in getting distributors " + err.Error(),
		})
	}

	count, err := dbconnection.Queries.GetDistributorsCount(dbconnection.Ctx, userId)
	if err != nil {
		log.Print(err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in getting distributors " + err.Error(),
		})
	}
	type Response struct {
		Data  []database.GetDistributorsRow `json:"data"`
		Count int64                         `json:"count"`
	}
	response := Response{
		Data:  data,
		Count: count,
	}
	// jsonString, err := json.Marshal(result)
	// if err != nil {
	// 	log.Printf("Error converting into Json: %s", err)
	// 	return c.Status(200).JSON(result)
	// }
	// if err := valkeyCache.SetValue("distributors:"+page+":"+userId, jsonString); err != nil {
	// 	log.Printf("Error in setting value in valkey %s ", err.Error())
	// }
	return c.Status(200).JSON(response)

}

func GetAllDistributors(c *fiber.Ctx) error {

	userId := c.Params("userId")

	data, err := dbconnection.Queries.GetAllDistributors(dbconnection.Ctx, userId)
	if err != nil {
		log.Print(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in getting Distributors " + err.Error(),
		})
	}

	// count, err := dbconnection.Queries.GetDistributorsCount(dbconnection.Ctx, userId)
	// if err != nil {
	// 	log.Print(err.Error())
	// 	return c.Status(500).JSON(fiber.Map{
	// 		"error": "Error in getting Distributors count " + err.Error(),
	// 	})
	// }
	// type Response struct {
	// 	Data  []database.GetAllDistributorsRow `json:"data"`
	// 	Count int64          	`json:"count"`
	// }
	// response := Response{
	// 	Data:  data,
	// 	Count: count,
	// }

	return c.Status(200).JSON(data)

}

func PostDistributor(c *fiber.Ctx) error {

	body := new(types.Distributor)
	userId := c.Params("userId")

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing body %s", err.Error())
		return c.Status(400).JSON("Error parssing body")
	}

	if err := dbconnection.Queries.PostDistributor(dbconnection.Ctx, database.PostDistributorParams{
		Name:     body.DistributorName,
		DomainID: body.DomainID,
		UserID:   userId,
	}); err != nil {
		log.Print(err.Error())
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return c.Status(409).JSON(fiber.Map{
					"error": "Distributor already exists, try another",
				})
			}
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in posting distributor " + err.Error(),
		})
	}

	// if err := valkeyCache.Revalidate("distributors:" + "1" + ":" + userId); err != nil {
	// 	log.Printf("Error in revalidating distributors cache: %s", err)
	// }
	return c.Status(201).JSON(fiber.Map{
		"msg": "distributor created successfully",
	})

}

func UpdateDistributor(c *fiber.Ctx) error {

	var params types.DistributorParams

	if err := c.ParamsParser(&params); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Error parsing distributorParams " + err.Error(),
		})
	}

	body := new(types.Distributor)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}

	if err := dbconnection.Queries.UpdateDistributor(dbconnection.Ctx, database.UpdateDistributorParams{
		ID:       params.DistributorId,
		UserID:   params.UserId,
		Name:     body.DistributorName,
		DomainID: body.DomainID,
	}); err != nil {
		log.Print(err.Error())
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return c.Status(409).JSON(fiber.Map{
					"error": "Distributor already exists, try another",
				})
			}
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in updating distributor " + err.Error(),
		})
	}
	// if err := valkeyCache.Revalidate("distributors:" + userId); err != nil {
	// 	log.Printf("Error in revalidating distributors cache: %s", err)
	// }
	return c.Status(200).JSON("distributor updated successfully")
}

func DeleteDistributor(c *fiber.Ctx) error {

	var params types.DistributorParams

	if err := c.ParamsParser(&params); err != nil {
		log.Print(err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Error in parsing parms :" + err.Error(),
		})
	}

	if err := dbconnection.Queries.DeleteDistributor(dbconnection.Ctx, database.DeleteDistributorParams{
		ID:     params.DistributorId,
		UserID: params.UserId,
	}); err != nil {
		var pgErr *pgconn.PgError
		log.Print(err.Error())
		if ok := errors.As(err, &pgErr); ok {
			switch code := pgErr.Code; code {
			case "23505":
				return c.Status(400).JSON(fiber.Map{
					"error": "No Record found",
				})
			case "23503":
				return c.Status(400).JSON(fiber.Map{
					"error": "Delete associated bills to delete distributor",
				})
			default:
				break
			}
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Error in deleting distributor " + err.Error(),
		})
	}

	// if err := valkeyCache.Revalidate("distributors:" + userId); err != nil {
	// 	log.Printf("Error in revalidating distributors cache: %s", err)
	// }
	return c.Status(200).JSON("Distributor deleted successfully")
}
