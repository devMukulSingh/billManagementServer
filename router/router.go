package router 

import (
	"github.com/devMukulSingh/billManagementServer.git/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1")

	v1.Post("/post-bill", controller.PostBill)
	v1.Post("/post-distributor", controller.PostDistributor)
	v1.Post("/post-domain", controller.PostDomain)

	v1.Put("/put-bill/:id",controller.UpdateBill)
	v1.Put("/put-domain/:id",controller.UpdateDomain)
	v1.Put("/put-distributor/:id",controller.UpdateDistributor)

	v1.Delete("/delete-domin/:id",controller.DeleteDomain)
	v1.Delete("/delete-distributor/:id",controller.DeleteDistributor)

	// v1.Post("/post-distributor", controller.PostDistributorController) 

}
