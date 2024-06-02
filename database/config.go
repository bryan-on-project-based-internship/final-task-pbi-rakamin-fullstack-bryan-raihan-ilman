package database

import (
	"fmt"
	"os"
	"user-profile-apis/app/models"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func Connect() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

    connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", dbHost, dbPort, dbUsername, dbPassword, dbName)
	fmt.Println(connectionString)

	db, err = gorm.Open("postgres", connectionString)

	db.AutoMigrate(&models.User{}, &models.Photo{})
		
	return db, err
}

// Close closes the database connection
func Close() error {
	sqlDB := db.DB()
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}
