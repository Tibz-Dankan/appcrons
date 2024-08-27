package tests

import (
	"log"
	"os"
	"testing"

	"github.com/Tibz-Dankan/keep-active/internal/events/publishers"
	"github.com/Tibz-Dankan/keep-active/internal/events/subscribers"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/schedulers"
	"github.com/Tibz-Dankan/keep-active/tests/setup"
)

type User = models.User
type App = models.User
type Request = models.User
type RequestTime = models.User
type Feedback = models.User

func TestMain(m *testing.M) {
	db := setup.DB
	err := db.AutoMigrate(&User{}, &App{}, &Request{}, &RequestTime{}, &Feedback{})
	if err != nil {
		log.Fatal("Failed to make auto migration", err)
	}
	log.Println("Auto Migration successful")

	go schedulers.InitSchedulers()
	go subscribers.InitEventSubscribers()
	publishers.InitEventPublishers()

	code := m.Run()

	setup.ClearAllTables()

	os.Exit(code)

	select {}
}
