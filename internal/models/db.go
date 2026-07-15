package models

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	gormDB          *gorm.DB
	once            sync.Once
	finishMigration sync.Once
)

func Db() *gorm.DB {
	once.Do(func() {
		var dsn string
		var err error

		env := os.Getenv("GO_ENV")
		log.Println("GO_ENV:", env)

		// Load dev .env file
		if env == "development" {
			err = godotenv.Load()
			if err != nil {
				log.Fatalf("Error loading .env file")
			}
			log.Println("Loaded .env var file")
		}

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

		// Phase 1: create/update every table and column, but suppress
		// creation of NEW foreign key constraints for this pass.
		// Location.userId, SiteVisit.userId and SiteVisit.locationId may
		// still contain blank ("") values written before the backfill in
		// FinishMigration runs (anonymous requests via
		// middlewares.OptionalAuth, or a failed geo-IP lookup) - with
		// constraints enabled, AutoMigrate would try (and fail) to ALTER
		// TABLE ... ADD CONSTRAINT against that stale data and log.Fatal
		// before the server ever starts.
		//
		// This must run over the full model list, not a subset that
		// excludes Location/SiteVisit: GORM's AutoMigrate auto-adds any
		// model a passed model depends on via FK even if it wasn't passed
		// (gorm.io/gorm migrator.ReorderModels, autoAdd=true) - Session
		// depends on Location, so excluding Location here would just get
		// it silently pulled back in with constraints still enabled.
		//
		// The backfill+constraint-creation pass (see FinishMigration)
		// can't happen here: it needs the package-level `db` var (used by
		// every models.* receiver method), which is itself initialized as
		// `var db = Db()` in schema.go - calling into those methods from
		// inside Db() would read `db` before its own initializer (this
		// function) has returned, an initialization cycle Go's compiler
		// rejects outright. FinishMigration must instead be called
		// explicitly once package initialization has completed (see
		// cmd/main.go).
		gormDB.DisableForeignKeyConstraintWhenMigrating = true
		err = gormDB.AutoMigrate(&User{}, &App{}, &Request{}, &RequestTime{},
			&Feedback{}, &OTP{}, &BugReport{}, &RequestCount{},
			&Location{}, &Session{}, &SiteVisit{})
		if err != nil {
			log.Fatal("Failed to make auto migration:", err)
		}
		log.Println("Auto Migration successful")

	})

	return gormDB
}

// FinishMigration completes the migration Db() intentionally left
// unfinished: it backfills any blank Location.userId/SiteVisit.userId
// (pointing them at the shared placeholder "unknown user" - see
// runUnknownUserBackfill) and blank SiteVisit.locationId, then re-runs
// AutoMigrate scoped to just Location/SiteVisit with foreign key
// constraint creation re-enabled, so the userId/locationId constraints
// that Db()'s pass deliberately skipped get created against clean data.
//
// Must be called once, explicitly, after package initialization has
// completed (see cmd/main.go) - it uses models.* receiver methods that
// depend on the package-level `db` var, which isn't safe to touch from
// inside Db() itself (see the comment there). Not re-running the full
// AutoMigrate model list here: any model whose constraint has never been
// created (e.g. BugReport.userId, which is always "" - see
// internal/routes/bugReport/postBugReport.go - out of scope for this fix)
// is left exactly as-is rather than newly surfacing a migration crash for
// a table this fix isn't meant to touch.
func FinishMigration() {
	finishMigration.Do(func() {
		if err := runUnknownUserBackfill(); err != nil {
			log.Fatal("Failed to backfill unknown user data:", err)
		}

		gormDB.DisableForeignKeyConstraintWhenMigrating = false
		if err := gormDB.AutoMigrate(&Location{}, &SiteVisit{}); err != nil {
			log.Fatal("Failed to make auto migration:", err)
		}
		log.Println("Auto Migration (foreign keys) successful")
	})
}
