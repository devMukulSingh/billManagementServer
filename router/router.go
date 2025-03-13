package router

import (
	"github.com/devMukulSingh/billManagementServer.git/controllers"
	"github.com/devMukulSingh/billManagementServer.git/middleware"
	"github.com/devMukulSingh/billManagementServer.git/types"

	// "github.com/devMukulSingh/billManagementServer.git/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")
	api.Post("/webhooks",controller.Webhook)
	//add user validation middleware
	v1 := api.Group("/v1/:userId")
	
	domain := v1.Group("/domain")
	bill := v1.Group("/bill")
	distributor := v1.Group("/distributor")
	product := v1.Group("/product")

	domain.Get("/get-all-domains",controller.GetAllDomains)
	domain.Get("/get-domains",middleware.ValidateQueryStrings, controller.GetDomains)
	domain.Post("/",middleware.ValidateBody[types.Domain](), controller.PostDomain)
	domain.Put("/:domainId",middleware.ValidateBody[types.Domain](),controller.UpdateDomain)
	domain.Delete("/:domainId",controller.DeleteDomain) 

	distributor.Get("/get-all-distributors",controller.GetAllDistributors)
	distributor.Get("/get-distributors",middleware.ValidateQueryStrings,controller.GetDistributors)
	distributor.Post("/", middleware.ValidateBody[types.Distributor](), controller.PostDistributor)
	distributor.Put("/:distributorId",middleware.ValidateBody[types.Distributor](),controller.UpdateDistributor)
	distributor.Delete("/:distributorId",controller.DeleteDistributor)

	// bill.Get("/get-all-bills",controller.GetAllBills)
	bill.Get("/get-bills",middleware.ValidateQueryStrings,controller.GetBills)
	bill.Post("/", middleware.ValidateBody[types.Bill](), controller.PostBill)
	bill.Put("/:billId", middleware.ValidateBody[types.Bill](), controller.UpdateBill)
	bill.Delete("/:billId",controller.DeleteBill)

	product.Get("/get-all-products",controller.GetAllProducts)
	product.Get("/get-products",middleware.ValidateQueryStrings,controller.GetProducts)
	product.Post("/",middleware.ValidateBody[types.Product](), controller.PostProduct)
	product.Put("/:productId",middleware.ValidateBody[types.Product](),controller.UpdateProduct)
	product.Delete("/:productId",controller.DeleteProduct)

}
