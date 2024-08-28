package config

import (
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func Db() *gorm.DB {
	once.Do(func() {
		var dsn string
		env := os.Getenv("GO_ENV")
		log.Println("GO_ENV:", env)
		var SkipDefaultTransaction bool = true

		switch env {
		case "development":
			dsn = os.Getenv("APPCRONS_DEV_DSN")
		case "testing":
			dsn = os.Getenv("APPCRONS_TEST_DSN")
			SkipDefaultTransaction = false
		case "staging":
			dsn = os.Getenv("APPCRONS_STAG_DSN")
			SkipDefaultTransaction = false
		case "production":
			dsn = os.Getenv("APPCRONS_PROD_DSN")
		default:
			log.Fatal("Unrecognized GO_ENV:", env)
		}

		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: SkipDefaultTransaction, PrepareStmt: true,
		})

		if err != nil {
			log.Fatal("Failed to connect to the database", err)
		}

		log.Println("Connected to database successfully")
	})

	return db
}
