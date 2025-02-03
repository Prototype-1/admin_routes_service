package utils

import (
	"fmt"
	"log"
	"github.com/Prototype-1/admin_routes_service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func InitDB() *gorm.DB {

	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}

	err = db.AutoMigrate(&models.Route{})
	if err != nil {
		log.Fatal("failed to auto-migrate Route model:", err)
	}

	log.Println("Database migrated successfully")

	return db
}
