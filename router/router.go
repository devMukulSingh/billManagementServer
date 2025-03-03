package router

import (
	"github.com/devMukulSingh/billManagementServer.git/controllers"
	"github.com/devMukulSingh/billManagementServer.git/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")
	api.Post("/webhooks",controller.Webhook)
	v1 := api.Group("/v1/:userId",middleware.ValidateUser)
	
	domain := v1.Group("/domain")
	bill := v1.Group("/bill")
	distributor := v1.Group("/distributor")
	product := v1.Group("/product")

	domain.Get("/get-all-domains",controller.GetAllDomains)
	domain.Get("/get-domains",controller.GetDomains)
	domain.Post("/post-domain", controller.PostDomain)
	domain.Put("/put-domain/:domainId",controller.UpdateDomain)
	domain.Delete("/delete-domain/:domainId",controller.DeleteDomain) 

	distributor.Get("/get-all-distributors",controller.GetAllDistributors)
	distributor.Get("/get-distributors",controller.GetDistributors)
	distributor.Post("/post-distributor", controller.PostDistributor)
	distributor.Put("/put-distributor/:distributorId",controller.UpdateDistributor)
	distributor.Delete("/delete-distributor/:distributorId",controller.DeleteDistributor)

	bill.Get("/get-all-bills",controller.GetAllBills)
	bill.Get("/get-bill",controller.GetBill)
	bill.Post("/post-bill", controller.PostBill)
	bill.Put("/put-bill/:billId",controller.UpdateBill)
	bill.Delete("/delete-bill/:billId",controller.DeleteBill)
	

	product.Get("/get-all-products",controller.GetAllProducts)
	product.Get("/get-products",controller.GetProducts)
	product.Post("/post-product",controller.PostProduct)
	product.Put("/put-product/:productId",controller.UpdateProduct)
	product.Delete("/delete-product/:productId",controller.DeleteProduct)

}
