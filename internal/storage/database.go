package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

type DatabaseInstance struct {
	Db *gorm.DB
}

var DB DatabaseInstance

// connectDb
func ConnectDb() {
	//make dsn from env

	dsn := os.Getenv("DSN")
	//log.Println("dsn: ", dsn)

	//dsn := "host=localhost user=postgres password=mypassabcbxghgwmmrmemrmdhdh dbname=go-db port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"
	//dsn := "host=localhost user=postgres password=mypass dbname=go-db port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}
	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")
	//db.AutoMigrate(&payment.PaymentModel{})

	DB = DatabaseInstance{
		Db: db,
	}
}

func Close() {
	sqlDB, err := DB.Db.DB()
	if err != nil {
		log.Printf("Error getting database instance: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
	log.Println("Database connection closed.")
}
