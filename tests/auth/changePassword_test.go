package tests

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
	"github.com/Tibz-Dankan/keep-active/tests/data"
	"github.com/Tibz-Dankan/keep-active/tests/setup"
)

func TestMissingPasswordsFields(t *testing.T) {
	time.Sleep(500 * time.Millisecond)
	setup.ClearAllTables()

	var label string = "Expects 400 with missing passwords"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	genData := data.NewGenTestData()
	name := genData.RandomUniqueName()
	email := genData.RandomUniqueEmail()
	password := genData.RandomUniquePassword(8)

	user := models.User{Name: name, Email: email, Password: password}

	userId, err := user.Create(user)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
		return
	}
	user.ID = userId

	signInPayload, _ := json.Marshal(user)
	req, _ = http.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(signInPayload))
	response = setup.ExecuteRequest(req)

	json.Unmarshal(response.Body.Bytes(), &body)

	token, ok := body["accessToken"]
	if !ok {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", token)
	}
	accessToken, _ := token.(string)
	bearerToken := fmt.Sprintf("Bearer %s", accessToken)

	path := fmt.Sprintf("/api/v1/auth/user/update-password/%s", userId)

	payload = []byte(`{"currentPassword":"password","newPassword":""}`)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

	payload = []byte(`{"currentPassword":"","newPassword":"newPassword"}`)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

	log.Printf("name:%s email:%s password:%s", name, email, password)
}

func TestSimilarPasswords(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 400 with current password similar to new password"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	genData := data.NewGenTestData()
	name := genData.RandomUniqueName()
	email := genData.RandomUniqueEmail()
	password := genData.RandomUniquePassword(8)

	user := models.User{Name: name, Email: email, Password: password}

	userId, err := user.Create(user)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new user. Got %v\n", err)
		return
	}

	signInPayload, _ := json.Marshal(user)
	req, _ = http.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(signInPayload))
	response = setup.ExecuteRequest(req)

	json.Unmarshal(response.Body.Bytes(), &body)

	token, ok := body["accessToken"]
	if !ok {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", token)
	}
	accessToken, _ := token.(string)
	bearerToken := fmt.Sprintf("Bearer %s", accessToken)

	path := fmt.Sprintf("/api/v1/auth/user/update-password/%s", userId)

	payload = []byte(`{"currentPassword":"password1","newPassword":"password1"}`)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)
	log.Printf("name:%s email:%s password:%s", name, email, password)
}

func TestNewPasswordSimilarToSavedOne(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 400 with new password similar to saved password"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	genData := data.NewGenTestData()
	name := genData.RandomUniqueName()
	email := genData.RandomUniqueEmail()
	password := genData.RandomUniquePassword(8)

	user := models.User{Name: name, Email: email, Password: password}

	userId, err := user.Create(user)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new user. Got %v\n", err)
		return
	}

	signInPayload, _ := json.Marshal(user)
	req, _ = http.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(signInPayload))
	response = setup.ExecuteRequest(req)

	json.Unmarshal(response.Body.Bytes(), &body)

	token, ok := body["accessToken"]
	if !ok {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", token)
	}
	accessToken, _ := token.(string)
	bearerToken := fmt.Sprintf("Bearer %s", accessToken)

	path := fmt.Sprintf("/api/v1/auth/user/update-password/%s", userId)

	payload = []byte(`{"currentPassword":"password","newPassword":"password"}`)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)
}

func TestSuccessfulChangePassword(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 200 on successful change password"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	genData := data.NewGenTestData()
	name := genData.RandomUniqueName()
	email := genData.RandomUniqueEmail()
	password := genData.RandomUniquePassword(8)

	user := models.User{Name: name, Email: email, Password: password}

	userId, err := user.Create(user)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new user. Got %v\n", err)
		return
	}

	signInPayload, _ := json.Marshal(user)
	req, _ = http.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(signInPayload))
	response = setup.ExecuteRequest(req)

	json.Unmarshal(response.Body.Bytes(), &body)

	token, ok := body["accessToken"]
	if !ok {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", token)
	}
	accessToken, _ := token.(string)
	bearerToken := fmt.Sprintf("Bearer %s", accessToken)

	path := fmt.Sprintf("/api/v1/auth/user/update-password/%s", userId)

	passwords := struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}{
		CurrentPassword: user.Password,
		NewPassword:     genData.RandomUniquePassword(8),
	}

	payload, _ = json.Marshal(passwords)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusOK, response.Code)
}
