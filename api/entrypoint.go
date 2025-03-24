package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)


var app *fiber.App

func init() {
    app = fiber.New()
	r := app.Group("/api");
    r.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello from Fiber on Vercel!")
    })
}

func Handler(w http.ResponseWriter, r *http.Request) {
    adaptor.FiberApp(app)(w, r)
}