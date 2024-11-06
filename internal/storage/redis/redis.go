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
			Password: rpassword, // пустая строка для без пароля
			DB:       rDB,
		})

		// Проверяем соединение с Redis
		_, err = rdb.Ping(ctx).Result()
		if err == nil {
			// Успешное подключение, выходим из цикла
			return rdb, nil
		}

		// Логируем ошибку и ждем перед повторной попыткой
		log.Printf("Ошибка подключения к Redis: %v. Попытка %d из %d", err, i+1, maxRetries)
		time.Sleep(retryInterval)
	}

	// Если попытки закончились и ошибка осталась, возвращаем ошибку
	return nil, fmt.Errorf("не удалось подключиться к Redis после %d попыток: %v", maxRetries, err)
}
