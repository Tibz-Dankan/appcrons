package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/tests/data"
	"github.com/Tibz-Dankan/keep-active/tests/setup"
)

func TestNonExistingForgotPasswordUser(t *testing.T) {
	setup.ClearAllTables()

	var label string
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder

	genData := data.NewGenTestData()

	email := struct {
		Email string `json:"email"`
	}{
		Email: genData.RandomUniqueEmail(),
	}

	label = "Expects 400 for non existing user"
	payload, _ = json.Marshal(email)
	req, _ = http.NewRequest("POST", "/api/v1/auth/forgot-password", bytes.NewBuffer(payload))
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)
}

func TestSuccessfulForgotPasswordInitialization(t *testing.T) {
	setup.ClearAllTables()

	var label string
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder

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
	time.Sleep(500 * time.Millisecond)

	emailStruct := struct {
		Email string `json:"email"`
	}{
		Email: email,
	}

	label = "Expects 200 on successful forgot password initialization"
	payload, _ = json.Marshal(emailStruct)
	req, _ = http.NewRequest("POST", "/api/v1/auth/forgot-password", bytes.NewBuffer(payload))
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusOK, response.Code)
}
