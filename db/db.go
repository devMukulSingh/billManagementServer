package database

import (
	"log"
	"os"
	
	"github.com/devMukulSingh/billManagementServer.git/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DbConn *gorm.DB

func ConnectDb() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	dsn := "host=" + host + " " + "user=" + user + " " + "password=" + password + " " + "dbname=" + dbName + " " + "port=" + port + " " + "sslmode=require"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
		TranslateError: true,
	})

	if err != nil {
		log.Fatal("Database connection failed")
	}
	if err := db.AutoMigrate(&model.User{},&model.Item{},&model.BillItem{},&model.Bill{},&model.Distributor{},&model.Domain{});err!= nil{
		log.Printf("Error migrating schema - %s",err.Error())
	}
	log.Print("Connection to db successfull")
	DbConn = db 
} 
