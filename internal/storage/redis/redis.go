package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const maxRetries = 10
const retryInterval = 5 * time.Second

// InitRedisCl - инициализирует клиент редиса
func InitRedisCl() (*redis.Client, error) {
	var rdb *redis.Client
	var err error
	ctx := context.Background()

	rhost := os.Getenv("APP_REDIS_HOST")
	rport := os.Getenv("APP_REDIS_PORT")
	rpassword := os.Getenv("APP_REDIS_PASSWORD")
	rDB, err := strconv.Atoi(os.Getenv("APP_REDIS_DB"))
	if err != nil {
		rDB = 0
	}

	for i := 0; i < maxRetries; i++ {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", rhost, rport),
			Password: rpassword,
			DB:       rDB,
		})

		_, err = rdb.Ping(ctx).Result()
		if err == nil {
			return rdb, nil
		}

		log.Printf("Ошибка подключения к Redis: %v. Попытка %d из %d", err, i+1, maxRetries)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("не удалось подключиться к Redis после %d попыток: %v", maxRetries, err)
}
