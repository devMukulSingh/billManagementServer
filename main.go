package main

import (
	"log"

	"github.com/devMukulSingh/billManagementServer.git/db"
	"github.com/devMukulSingh/billManagementServer.git/lib"
	"github.com/devMukulSingh/billManagementServer.git/router"
	"github.com/devMukulSingh/billManagementServer.git/valkeyCache"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading dotenv")
	}

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: constants.BASE_URL_CLIENT,
	}))

	if err := valkeyCache.Connect();err!=nil{
		log.Printf("Error connecting to valkey : %s",err.Error())
	}

	database.ConnectDb()

	db, err := database.DbConn.DB()

	if err != nil {
		log.Fatal("Error in DB")
	}

	defer db.Close()

	router.SetupRoutes(app)

	log.Print("Server is running at 8000")
	app.Listen(":8000")
}
