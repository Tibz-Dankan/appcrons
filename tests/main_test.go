package tests

import (
	"os"
	"testing"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/tests/setup"
)

func TestMain(m *testing.M) {
	models.DBAutoMigrate()

	code := m.Run()

	setup.ClearAllTables()

	os.Exit(code)
}
