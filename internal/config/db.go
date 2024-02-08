package config

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Db() *gorm.DB {

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	dsn := os.Getenv("DSN_KEEP_ACTIVE")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true, PrepareStmt: true,
	})

	if err != nil {
		log.Fatal("Failed to connect to the database", err)
	}

	fmt.Println("Connected to database successfully")
	return db
}

// gorm migrate create -name=create-users
