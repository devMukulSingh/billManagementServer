package api

import (
	"net/http"
	"log"
	 "github.com/devMukulSingh/billManagementServer.git/dbConnection"
	// "github.com/devMukulSingh/billManagementServer.git/lib"
	"github.com/devMukulSingh/billManagementServer.git/router"
	// "github.com/devMukulSingh/billManagementServer.git/valkeyCache"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

var app *fiber.App

func init() {

	app = fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	// if err := valkeyCache.Connect();err!=nil{
	// 	log.Printf("Error connecting to valkey : %s",err.Error())
	// }

	if err := dbconnection.ConnectDb(); err!=nil{
		log.Fatalf("Error in connection db : %s",err.Error())
	}

	log.Print("Db connection successfull")
	router.SetupRoutes(app)
	// log.Print("Server is running at 8000")
	app.Listen(":8000")
}

func Handler(w http.ResponseWriter, r *http.Request) {
    adaptor.FiberApp(app)(w, r)
	// defer dbconnection.Connection.Close()
}