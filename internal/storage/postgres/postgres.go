package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/Mallbrusss/BackEntryMiddle/models"

	"github.com/Mallbrusss/BackEntryMiddle/pkg/envloader"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


const envPathToStorageEnv string = "../../deployment/.env"

func InitDB() *gorm.DB {
	envloader.LoadEnv(envPathToStorageEnv)
	log.Println("Loaded PSQL storage.env file")

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dBname := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dBname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to db: %v", err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Document{})
	db.AutoMigrate(&models.DocumentAccess{})

	return db
}
