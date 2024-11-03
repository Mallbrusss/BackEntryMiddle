package envbalancer

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv загружает переменные окружения из указанного файла
func LoadEnv(filename string) {
	err := godotenv.Load(filename)
	if err != nil {
		log.Fatalf("Error loading %s file: %v", filename, err)
	}
}

// GetEnv получает значение переменной окружения по ключу
func GetEnv(key string) string {
	return os.Getenv(key)
}
