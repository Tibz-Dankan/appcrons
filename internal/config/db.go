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

	dsn := os.Getenv("DSN_RESERVE_NOW")

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
