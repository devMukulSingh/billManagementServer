package router

import (
	"github.com/devMukulSingh/billManagementServer.git/controllers"
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
	domain.Get("/get-domains",controller.GetDomains)
	domain.Post("/", controller.PostDomain)
	domain.Put("/:domainId",controller.UpdateDomain)
	domain.Delete("/:domainId",controller.DeleteDomain) 

	distributor.Get("/get-all-distributors",controller.GetAllDistributors)
	distributor.Get("/get-distributors",controller.GetDistributors)
	distributor.Post("/", controller.PostDistributor)
	distributor.Put("/:distributorId",controller.UpdateDistributor)
	distributor.Delete("/:distributorId",controller.DeleteDistributor)

	// bill.Get("/get-all-bills",controller.GetAllBills)
	bill.Get("/get-bills",controller.GetBills)
	bill.Post("/", controller.PostBill)
	bill.Put("/:billId",controller.UpdateBill)
	bill.Delete("/:billId",controller.DeleteBill)

	product.Get("/get-all-products",controller.GetAllProducts)
	product.Get("/get-products",controller.GetProducts)
	product.Post("/",controller.PostProduct)
	product.Put("/:productId",controller.UpdateProduct)
	product.Delete("/:productId",controller.DeleteProduct)

}
