package config

import (
	"fmt"
	"log"
	"os"

	"github.com/golang/glog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"backEntryMiddle/envBalancer"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
}

func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{}
}

// путь до env файла указывается относительно исполняемого файла
// если вызываем из tinkoff_test, то указываем относительный путь tinkoff_test
// TODO: перенести вызыв env в корень проекта

const envPathToStorageEnv string = "../../deployment/.env"

func (pg *PostgresConfig) InitDB() *gorm.DB {
	envbalancer.LoadEnv(envPathToStorageEnv)
	glog.Infof("Loaded storage.env file")

	pg.Host = os.Getenv("POSTGRES_HOST")
	pg.Port = os.Getenv("POSTGRES_PORT")
	pg.User = os.Getenv("POSTGRES_USER")
	pg.Password = os.Getenv("POSTGRES_PASSWORD")
	pg.DBname = os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", pg.Host, pg.User, pg.Password, pg.DBname, pg.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to db: %v", err)
	}
	return db
}
