package setup

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tibz-Dankan/keep-active/internal/config"
	"github.com/Tibz-Dankan/keep-active/internal/routes"

	_ "github.com/lib/pq"
)

var db = config.Db()
var DB = db

func ClearAllTables() {
	tables := []string{"requests", "request_times", "apps", "feedbacks", "users"}

	for _, table := range tables {
		statement := fmt.Sprintf("DELETE FROM %s ;", table)
		db.Exec(statement)
	}
}

func ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r := routes.AppRouter()
	r.ServeHTTP(rr, req)

	return rr
}

func CheckResponseCode(t *testing.T, label string, expected, actual int) {
	if expected != actual {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	} else {
		fmt.Printf("=== PASS: %s\n", label)
	}
}
