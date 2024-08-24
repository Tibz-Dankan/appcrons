package config

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Db() *gorm.DB {
	var dsn string

	env := os.Getenv("GO_ENV")
	log.Println("GO_ENV:", env)

	if env == "development" {
		dsn = os.Getenv("APPCRONS_DEV_DSN")

	} else if env == "testing" {
		dsn = os.Getenv("APPCRONS_TEST_DSN")

	} else if env == "staging" {
		dsn = os.Getenv("APPCRONS_STAG_DSN")

	} else if env == "production" {
		dsn = os.Getenv("APPCRONS_PROD_DSN")
	}

	log.Println("dsn:", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true, PrepareStmt: true,
	})

	if err != nil {
		log.Fatal("Failed to connect to the database", err)
	}

	log.Println("Connected to database successfully")
	return db
}
