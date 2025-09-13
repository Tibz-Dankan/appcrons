package models

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func RedisClient() *redis.Client {
    var REDIS_URL string

	env := os.Getenv("GO_ENV")
	if env == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
		log.Println("Loaded .env var file")
	}

	switch env {
	case "development":
		REDIS_URL = os.Getenv("REDIS_DEV_URL")
	case "production":
		REDIS_URL = os.Getenv("REDIS_PROD_URL")
	default:
		log.Fatal("Unrecognized GO_ENV:", env)
	}


    log.Println("REDIS_URL:", REDIS_URL)

	opt, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		log.Fatal("Failed to connect to redis", err)
	}

	client := redis.NewClient(opt)
	log.Println("Connected to redis successfully")

	return client
}
