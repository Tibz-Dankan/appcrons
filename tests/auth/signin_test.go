package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/tests/data"
	"github.com/Tibz-Dankan/keep-active/tests/setup"
)

func TestInvalidSignInCredentials(t *testing.T) {
	setup.ClearAllTables()

	var label string
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder

	label = "Expects 400 with missing/invalid email"
	payload = []byte(`{"email":"","password":"password"}`)
	req, _ = http.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(payload))
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

	label = "Expects 400 with missing/invalid password"
	payload = []byte(`{"email":"user@gmail.com","password":""}`)
	req, _ = http.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(payload))
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)
}

func TestSuccessfulSignIn(t *testing.T) {
	setup.ClearAllTables()

	var label string
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	genData := data.NewGenTestData()
	name := genData.RandomUniqueName()
	email := genData.RandomUniqueEmail()
	password := genData.RandomUniquePassword(8)

	user := models.User{Name: name, Email: email, Password: password}

	_, err := user.Create(user)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new user. Got %v\n", err)
		return
	}

	label = "Expects accessToken on successful signin"
	payload, _ = json.Marshal(user)
	req, _ = http.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(payload))
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusOK, response.Code)

	json.Unmarshal(response.Body.Bytes(), &body)

	token, found := body["accessToken"]
	if !found {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", token)
	}

	accessToken, ok := token.(string)
	if !ok {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", accessToken)
	}
	if accessToken == "" {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got an empty string")
	} else {
		fmt.Printf("=== PASS: %s\n", label)
	}
}
