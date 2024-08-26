package tests

import (
	"os"
	"testing"

	"github.com/Tibz-Dankan/keep-active/internal/events/publishers"
	"github.com/Tibz-Dankan/keep-active/internal/events/subscribers"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/schedulers"
	"github.com/Tibz-Dankan/keep-active/tests/setup"
)

func TestMain(m *testing.M) {
	models.DBAutoMigrate()

	go schedulers.InitSchedulers()
	go subscribers.InitEventSubscribers()
	publishers.InitEventPublishers()

	code := m.Run()

	setup.ClearAllTables()

	os.Exit(code)

	select {}
}
