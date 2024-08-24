package tests

import (
	"os"
	"testing"

	"github.com/Tibz-Dankan/keep-active/internal/models"
)

func TestMain(m *testing.M) {
	models.DBAutoMigrate()

	code := m.Run()

	models.DBDropTables()

	os.Exit(code)
}
