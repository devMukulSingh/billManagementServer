package middleware

import (
	"github.com/devMukulSingh/billManagementServer.git/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateUser(c *fiber.Ctx) error {
	var params types.Param
	if err := c.ParamsParser(&params); err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error": "Invalid userId",
		})
	}
	return c.Next()
}

var validate = validator.New()

func GetValidationErrors(validationStruct interface{}) []types.IError {
	var validationErrors = []types.IError{}
	if errs := validate.Struct(validationStruct); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem types.IError
			elem.Field = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Param()
			validationErrors = append(validationErrors, elem)
		}
		return validationErrors
	}
	return nil
}

func ValidateQueryStrings(c *fiber.Ctx) error {
	query := new(types.Query)
	c.QueryParser(query)
	if validationErrors := GetValidationErrors(query); validationErrors != nil {
		return c.Status(400).JSON(validationErrors)
	}
	return c.Next()
}

func ValidateBody[BodyType any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(BodyType)
		c.BodyParser(body)
		if validationErrors := GetValidationErrors(body); validationErrors != nil {
			return c.Status(400).JSON(validationErrors)
		}
		return c.Next()
	}
}

// func ValidatePostDomain(c *fiber.Ctx) error {
// 	body := new(types.Domain)
// 	c.BodyParser(body)
// 	if validationErrors := GetValidationErrors(body); validationErrors != nil {
// 		return c.Status(400).JSON(validationErrors)
// 	}
// 	return c.Next()
// }

// func ValidatePostBill(c *fiber.Ctx) error {
// 	body := new(types.Bill)
// 	c.BodyParser(body)
// 	if validationErrors := GetValidationErrors(body); validationErrors != nil {
// 		return c.Status(400).JSON(validationErrors)
// 	}
// 	return c.Next()
// }

// func ValidatePostDistributor(c *fiber.Ctx) error {

// 	body := new(types.Distributor)
// 	c.BodyParser(body)
// 	if validationErrors := GetValidationErrors(body); validationErrors != nil {
// 		return c.Status(400).JSON(validationErrors)
// 	}
// 	return c.Next()
// }

// func ValidatePostProduct(c *fiber.Ctx) error {

// 	body := new(types.Product)
// 	c.BodyParser(body)
// 	if validationErrors := GetValidationErrors(body); validationErrors != nil {
// 		return c.Status(400).JSON(validationErrors)
// 	}
// 	return c.Next()
// }
