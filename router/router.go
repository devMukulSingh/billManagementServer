package router

import (
	"github.com/devMukulSingh/billManagementServer.git/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/",controller.PostBillController)
}