package controller

import (
	"database/sql"
	"errors"
	"log"

	"github.com/devMukulSingh/billManagementServer.git/database"
	"github.com/devMukulSingh/billManagementServer.git/dbConnection"
	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/jackc/pgx/v5/pgconn"

	// "github.com/devMukulSingh/billManagementServer.git/valkeyCache"
	"github.com/gofiber/fiber/v2"
)

func GetAllDomains(c *fiber.Ctx) error {

	var params types.DomainParams
	if err := c.ParamsParser(&params); err!=nil{
		log.Printf("Error parsing userId in params : %s",err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error":"Error parsing userId in params",
		})
	}

	type Response struct {
		Data  interface{} 		`json:"data"`
		Count int64          	`json:"count"`
	}
	data,err := dbconnection.Queries.GetAllDomains(dbconnection.Ctx,params.UserID);

	if err!=nil{
		log.Print(err.Error())
	}

	response := Response{
		Data:  data.Data,
		Count: data.Count,
	}

	return c.Status(200).JSON(response)
}
func GetDomains(c *fiber.Ctx) error {

	userId := c.Params("userId")
	page := int32(c.QueryInt("page", 1))
	limit := int32(c.QueryInt("limit", 10))

	if err := c.QueryParser(&types.Query{}); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	//valkey cache
	// cache, err := valkeyCache.GetValue("domains:" + userId)
	// if err != nil {
	// 	if err.Error() != "valkey nil message" {
	// 		log.Printf("Error in getting cached bills : %s", err)
	// 	}
	// } else {
	// 	c.Set("Content-Type", "application/json")
	// 	return c.Status(200).SendString(cache)
	// }

	data,err := dbconnection.Queries.GetDomains(dbconnection.Ctx,database.GetDomainsParams{
		UserID: userId,
		Offset: (page - 1) * limit,
		Limit: limit,
	})
	
	if err!=nil{
		log.Fatalf("Error getting domains : %s", err.Error());
	}
	//TODO: optmimise more to get count in a single query
	count,err := dbconnection.Queries.GetDomainsCount(dbconnection.Ctx,userId)
	if err!=nil{
		log.Fatalf("Error getting domains : %s", err.Error());
	}
	
	// type Data struct{
	// 	Name		string			`json:"name"`
	// 	Id 			pgtype.UUID		`json:"id"`
	// 	CreatedAt	pgtype.UUID		`json:"created_at"`
	// }
	type Response struct {
		Data  []database.GetDomainsRow			 `json:"data"`
		Count int64          					`json:"count"`
	}
	response := Response{
		Data:  data,
		Count: count,
	}
	// jsonDomain, err := json.Marshal(domains)
	// if err != nil {
	// 	log.Print("error converting to json")
	// }
	// if err := valkeyCache.SetValue("domains:"+userId, jsonDomain); err != nil {
	// 	log.Printf("Error in setting value in valkey %s ", err.Error())
	// }
	return c.Status(200).JSON(response)
}

func PostDomain(c *fiber.Ctx) error {

	params := struct{ UserId    string	`params:"userId"`}{}
	if err:= c.ParamsParser(&params); err!= nil{
		log.Print(err)
		return c.Status(400).JSON(fiber.Map{
			"error":"Error parsing params " + err.Error(),
		})
	}
	body := new(types.Domain)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}
	
	if err := dbconnection.Queries.PostDomain(dbconnection.Ctx,database.PostDomainParams{
		Name: body.DomainName,
		UserID: params.UserId,
	}); err!=nil{
		var pgErr *pgconn.PgError
		log.Print(err.Error())
		if ok := errors.As(err,&pgErr) ; ok{
			if pgErr.Code=="23505"{
				return c.Status(400).JSON(fiber.Map{
					"error":"Domain already exists, try another.",
				})
			}
		}
		return c.Status(500).JSON(fiber.Map{
			"error":"Error posting domain : " + err.Error(),
		})
	}

	// result := database.DbConn.Create(&model.Domain{
	// 	Name:   body.DomainName,
	// 	UserID: userId,
	// })
	// if result.Error != nil {
	// 	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
	// 		log.Print(result.Error.Error())
	// 		return c.Status(409).JSON(fiber.Map{
	// 			"error": "Domain already exists, try another",
	// 		})
	// 	}
	// 	log.Printf("Error in saving Domain into db %s", result.Error.Error())
	// 	return c.Status(500).JSON("Internal server error")
	// }

	//revalidate cache
	// valkeyCache.Revalidate("domains:" + userId)
	return c.Status(201).JSON(fiber.Map{
		"msg": "domain created successfully",
	})
}

func UpdateDomain(c *fiber.Ctx) error {
	var params types.DomainParams;

	if err := c.ParamsParser(&params);err!=nil{
		return c.Status(400).JSON(fiber.Map{
			"error":"Error parsing params :" + err.Error(),
		})
	}

	body := new(types.Domain)

	if err := c.BodyParser(body); err != nil {
		log.Printf("Error parsing req body %s", err.Error())
		return c.Status(400).JSON("Error parsing body")
	}

	if err := dbconnection.Queries.UpdateDomain(dbconnection.Ctx,database.UpdateDomainParams{
		ID: params.DomainID,
		UserID: params.UserID,
		Name: body.DomainName,
	}); err!=nil{
		var pgErr *pgconn.PgError
		log.Print(err.Error())
		if ok := errors.As(err,&pgErr) ; ok{
			if pgErr.Code=="23505"{
				return c.Status(400).JSON(fiber.Map{
					"error":"Domain already exists, try another.",
				})
			}
		}
		return c.Status(500).JSON(fiber.Map{
			"error":"Error updating domain : " + err.Error(),
		})
	}

	// if result := database.DbConn.Model(&model.Domain{}).Where("id=? AND user_id=?", domainId, userId).Update("name", body.DomainName); result.Error != nil {
	// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 		log.Printf("No domain found %s", result.Error.Error())
	// 		return c.Status(400).JSON("No domain found")
	// 	}
	// 	log.Printf("Error updating domain %s", result.Error.Error())
	// 	return c.Status(500).JSON("Error updating domain")
	// }
	// if err := valkeyCache.Revalidate("domains:" + userId); err != nil {
	// 	log.Printf("Error in revalidating domains cache: %s", err)
	// }
	return c.Status(200).JSON("domain updated successfully")
}

func DeleteDomain(c *fiber.Ctx) error {
	
	var params types.DomainParams;

	if err := c.ParamsParser(&params);err!=nil{
		log.Print(err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error":"Error parsing params " + err.Error(),
		})
	}

	if err := dbconnection.Queries.DeleteDomain(dbconnection.Ctx,database.DeleteDomainParams{
		ID: params.DomainID,
		UserID: params.UserID,
	}); err!=nil{
		log.Print(err.Error())
		if errors.Is(err,sql.ErrNoRows){
			return c.Status(404).JSON(fiber.Map{
				"error":"no domain found ",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error":"Error deleting domain : " + err.Error(),
		})
	}
	// if err := database.DbConn.Where("id =? AND user_id=?", domainId, userId).Delete(&model.Domain{}).Error; err != nil {
	// 	log.Printf("Error deleting domain %s", err.Error())

	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		log.Printf("No domain found %s", err.Error())
	// 		return c.Status(400).JSON("Error:No record found")
	// 	}
	// 	if errors.Is(err, gorm.ErrForeignKeyViolated) {
	// 		return c.Status(405).JSON(fiber.Map{
	// 			"error": "Delete associated bills and distributors to delete domain",
	// 		})
	// 	}

	// 	return c.Status(500).JSON("Error deleting domain")
	// }
	// if err := valkeyCache.Revalidate("domains:" + userId); err != nil {
	// 	log.Printf("Error in revalidating domains cache: %s", err)
	// }
	return c.Status(200).JSON("domain deleted successfully")
}

