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

	domain.Get("/get-domain",controller.GetDomain)
	domain.Post("/post-domain", controller.PostDomain)
	domain.Put("/put-domain/:id",controller.UpdateDomain)
	domain.Delete("/delete-domain/:id",controller.DeleteDomain)

	distributor.Post("/post-distributor", controller.PostDistributor)
	distributor.Put("/put-distributor/:id",controller.UpdateDistributor)
	distributor.Delete("/delete-distributor/:id",controller.DeleteDistributor)

	bill.Post("/post-bill", controller.PostBill)
	bill.Put("/put-bill/:id",controller.UpdateBill)
	bill.Delete("/delete-bill/:id",controller.DeleteBill)
	
	item.Post("/post-item",controller.PostItem)
	item.Put("/put-item/:id",controller.UpdateItem)
	item.Delete("/delete-item/:id",controller.DeleteItem)

}
