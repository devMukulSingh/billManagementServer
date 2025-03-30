package main

import (
	"log"
	// "github.com/devMukulSingh/billManagementServer.git/db"
	 "github.com/devMukulSingh/billManagementServer.git/dbConnection"
	// "github.com/devMukulSingh/billManagementServer.git/lib"
	"github.com/devMukulSingh/billManagementServer.git/router"
	"github.com/devMukulSingh/billManagementServer.git/valkeyCache"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	// "gorm.io/gorm"
)
var app * fiber.App;
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("Error loading dotenv : " + err.Error())
	}

	app = fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	if err := valkeyCache.Connect();err!=nil{
		log.Printf("Error connecting to valkey : %s",err.Error())
	}

	if err := dbconnection.ConnectDb(); err!=nil{
		log.Fatalf("Error in connection db : %s",err.Error())
	}
	defer dbconnection.Connection.Close()
	
	router.SetupRoutes(app)
	log.Print("Server is running at 8000")
	app.Listen(":8000")
}
