package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func RedisClient() *redis.Client {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	REDIS_URL := os.Getenv("REDIS_URL")

	fmt.Println("REDIS_URL==>", REDIS_URL)

	opt, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		log.Fatal("Failed to connect to redis", err)
	}

	client := redis.NewClient(opt)

	return client
}
