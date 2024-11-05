package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Mallbrusss/BackEntryMiddle/pkg/envloader"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

const envPathToStorageEnv string = "../../deployment/.env"

func InitRedisCl() *redis.Client {
	envloader.LoadEnv(envPathToStorageEnv)
	log.Println("Loaded Redis storage.env file")

	rhost := os.Getenv("REDIS_HOST")
	rport := os.Getenv("REDIS_PORT")
	// ruser := os.Getenv("REDIS_USER")
	rpassword := os.Getenv("REDIS_PASSWORD")

	rDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		rDB = 0
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rhost, rport),
		Password: rpassword,
		DB:       rDB,
	})

	_, errR := rdb.Ping(ctx).Result()
	if errR != nil {
		// log.Println("error conect to redis:", errR)
		// return nil

		log.Fatal(errR)
	}

	log.Println("Success connect to Redis")
	return rdb
}
