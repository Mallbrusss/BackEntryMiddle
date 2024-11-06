package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Mallbrusss/BackEntryMiddle/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const maxRetries = 10
const retryInterval = 5 * time.Second

func InitDB() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	host := os.Getenv("APP_POSTGRES_HOST")
	port := os.Getenv("APP_POSTGRES_PORT")
	user := os.Getenv("APP_POSTGRES_USER")
	password := os.Getenv("APP_POSTGRES_PASSWORD")
	dBname := os.Getenv("APP_POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dBname, port)

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		log.Printf("Ошибка подключения к базе данных: %v. Попытка %d из %d", err, i+1, maxRetries)
		time.Sleep(retryInterval)
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

	return db, err
}
