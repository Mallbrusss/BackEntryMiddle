package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/Mallbrusss/BackEntryMiddle/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// const envPathToStorageEnv string = "./../../.env"

func InitDB() *gorm.DB {
	host := os.Getenv("APP_POSTGRES_HOST")
	port := os.Getenv("APP_POSTGRES_PORT")
	user := os.Getenv("APP_POSTGRES_USER")
	password := os.Getenv("APP_POSTGRES_PASSWORD")
	dBname := os.Getenv("APP_POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dBname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to db: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Println("error migrate user table")
	}

	err = db.AutoMigrate(&models.Document{})
	if err != nil {
		log.Println("error migrate Document table")
	}

	err = db.AutoMigrate(&models.DocumentAccess{})
	if err != nil {
		log.Println("error migrate DocumentAccess table")
	}

	log.Println("Success connect to Postgres")
	return db
}
