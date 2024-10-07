package models

import (
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	gormDB *gorm.DB
	once   sync.Once
)

func Db() *gorm.DB {
	once.Do(func() {
		var dsn string
		var err error

		env := os.Getenv("GO_ENV")
		log.Println("GO_ENV:", env)

		switch env {
		case "development":
			dsn = os.Getenv("APPCRONS_DEV_DSN")
		case "testing":
			dsn = os.Getenv("APPCRONS_TEST_DSN")
		case "staging":
			dsn = os.Getenv("APPCRONS_STAG_DSN")
		case "production":
			dsn = os.Getenv("APPCRONS_PROD_DSN")
		default:
			log.Fatal("Unrecognized GO_ENV:", env)
		}

		isTestingOrStagingEnv := env == "testing" || env == "staging"

		// Use sqlite db for testing
		if isTestingOrStagingEnv {
			gormDB, err = gorm.Open(sqlite.Open("./../../appcrons_test.db"), &gorm.Config{
				SkipDefaultTransaction: true, PrepareStmt: true,
				// Logger: logger.Default.LogMode(logger.Info),
			})
			if err != nil {
				log.Fatal("Failed to connect to sqlite db:", err)
			}
			log.Println("Connected to sqlite successfully")

		}

		// Use postgres as the primary db in dev and prod
		if !isTestingOrStagingEnv {
			gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				SkipDefaultTransaction: true, PrepareStmt: true,
			})

			if err != nil {
				log.Fatal("Failed to connect to the database", err)
			}
			log.Println("Connected to postgres successfully")

		}

		err = gormDB.AutoMigrate(&User{}, &App{}, &Request{}, &RequestTime{}, &Feedback{}, &OTP{}, &BugReport{})
		if err != nil {
			log.Fatal("Failed to make auto migration", err)
		}
		log.Println("Auto Migration successful")

	})

	return gormDB
}
