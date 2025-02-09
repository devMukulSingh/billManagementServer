package router

import (
	"github.com/devMukulSingh/billManagementServer.git/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1")

	v1.Post("/post-bill", controller.PostBillController)
	v1.Post("/post-distributor", controller.PostDistributorController)
	v1.Post("/post-domain", controller.PostDomainController)
	// v1.Post("/post-distributor", controller.PostDistributorController)

}
