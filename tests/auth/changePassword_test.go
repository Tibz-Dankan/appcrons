package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/tests/setup"
)

func TestMissingPasswordsFields(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 400 with missing passwords"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	user := models.User{Name: "username", Email: "user@gmail.com", Password: "password"}

	userId, err := user.Create(user)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
		return
	}

	signInPayload := []byte(`{"email":"user@gmail.com","password":"password"}`)
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
}

func TestSimilarPasswords(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 400 with current password similar to new password"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	user := models.User{Name: "username", Email: "user@gmail.com", Password: "password"}

	userId, err := user.Create(user)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
		return
	}

	signInPayload := []byte(`{"email":"user@gmail.com","password":"password"}`)
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
}

func TestNewPasswordSimilarToSavedOne(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 400 with new password similar to saved password"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	user := models.User{Name: "username", Email: "user@gmail.com", Password: "password"}

	userId, err := user.Create(user)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
		return
	}

	signInPayload := []byte(`{"email":"user@gmail.com","password":"password"}`)
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

	user := models.User{Name: "username", Email: "user@gmail.com", Password: "password"}

	userId, err := user.Create(user)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
		return
	}

	signInPayload := []byte(`{"email":"user@gmail.com","password":"password"}`)
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

	payload = []byte(`{"currentPassword":"password","newPassword":"newPassword"}`)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusOK, response.Code)
}
