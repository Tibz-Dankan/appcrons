package config

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func RedisClient() *redis.Client {

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	REDIS_URL := os.Getenv("REDIS_URL")

	opt, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		log.Fatal("Failed to connect to redis", err)
	}

	client := redis.NewClient(opt)
	log.Println("Connected to redis successfully")

	return client
}
