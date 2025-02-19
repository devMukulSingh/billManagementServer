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
	item := v1.Group("/item")

	domain.Get("/get-all-domains",controller.GetAllDomains)
	domain.Get("/get-domain/:domainId",controller.GetDomain)
	domain.Post("/post-domain", controller.PostDomain)
	domain.Put("/put-domain/:domainId",controller.UpdateDomain)
	domain.Delete("/delete-domain/:domainId",controller.DeleteDomain)

	distributor.Get("/get-all-distributors",controller.GetAllDistributors)
	distributor.Get("/get-distributor/:distributorId",controller.GetDistributor)
	distributor.Post("/post-distributor", controller.PostDistributor)
	distributor.Put("/put-distributor/:distributorId",controller.UpdateDistributor)
	distributor.Delete("/delete-distributor/:distributorId",controller.DeleteDistributor)

	bill.Get("/get-all-bills",controller.GetAllBills)
	bill.Get("/get-bill",controller.GetBill)
	bill.Post("/post-bill", controller.PostBill)
	bill.Put("/put-bill/:billId",controller.UpdateBill)
	bill.Delete("/delete-bill/:billId",controller.DeleteBill)
	

	item.Get("/get-all-items/:billId",controller.GetAllItems)
	// item.Get("/get-item",controller.GetItem)
	item.Post("/post-item",controller.PostItem)
	item.Put("/put-item/:itemId",controller.UpdateItem)
	item.Delete("/delete-item/:itemId",controller.DeleteItem)

}
