package setup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/routes"
	"github.com/Tibz-Dankan/keep-active/tests/data"

	_ "github.com/lib/pq"
)

var db = models.DB

// Deletes data for all tables in the application
func ClearAllTables() {
	tables := []string{"requests", "request_times", "apps", "feedbacks", "users"}

	for _, table := range tables {
		statement := fmt.Sprintf("DELETE FROM %s;", table)
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

// This function returns 'userId, email, password, bearerToken' and should
// be used by test functions that require authentication
// details like userId and bearerToken
func CreateSignInUser() (string, string, string, string) {
	var label string = "Sign in user"
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	genData := data.NewGenTestData()
	name := genData.RandomUniqueName()
	email := genData.RandomUniqueEmail()
	password := genData.RandomUniquePassword(15)

	user := models.User{Name: name, Email: email, Password: password}

	userId, err := user.Create(user)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		errorMessage := fmt.Errorf("expects new user. Got %v", err)
		fmt.Println(errorMessage)
	}
	time.Sleep(500 * time.Millisecond)

	signInStruct := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	signInPayload, _ := json.Marshal(signInStruct)
	req, _ = http.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(signInPayload))
	response = ExecuteRequest(req)

	json.Unmarshal(response.Body.Bytes(), &body)

	log.Println("Body: ", body)

	token, ok := body["accessToken"]
	if !ok {
		fmt.Printf("=== FAIL: %s\n", label)
		errorMessage := fmt.Errorf("expects accessToken. Got %v", token)
		fmt.Println(errorMessage)
	}
	accessToken, _ := token.(string)
	bearerToken := fmt.Sprintf("Bearer %s", accessToken)

	return userId, email, password, bearerToken
}

// This function returns 'userId, bearerToken'and should
// be used by test functions that require authentication
// details like userId and bearerToken
func SignInUser(email, password string) (string, string) {
	var label string = "Sign in user"
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	signInStruct := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	signInPayload, _ := json.Marshal(signInStruct)
	req, _ = http.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(signInPayload))
	response = ExecuteRequest(req)

	json.Unmarshal(response.Body.Bytes(), &body)

	log.Println("Body: ", body)

	token, ok := body["accessToken"]
	if !ok {
		fmt.Printf("=== FAIL: %s\n", label)
		errorMessage := fmt.Errorf("expects accessToken. Got %v", token)
		fmt.Println(errorMessage)
	}
	accessToken, _ := token.(string)
	bearerToken := fmt.Sprintf("Bearer %s", accessToken)

	user, ok := body["user"].(map[string]interface{})
	if !ok {
		fmt.Printf("=== FAIL: %s\n", label)
		errorMessage := fmt.Errorf("expects user. Got %v", user)
		fmt.Println(errorMessage)

	}

	userId, ok := user["id"].(string)
	if !ok {
		fmt.Printf("=== FAIL: %s\n", label)
		errorMessage := fmt.Errorf("expects user id. Got %v", userId)
		fmt.Println(errorMessage)
	}

	return userId, bearerToken
}
