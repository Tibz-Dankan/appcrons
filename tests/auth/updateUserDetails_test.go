package tests

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

func TestMissingUpdateUserDetailFields(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 400 with missing name/email"
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

	path := fmt.Sprintf("/api/v1/auth/user/update/%s", userId)

	payload = []byte(`{"name":"username","email":""}`)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

	payload = []byte(`{"name":"","email":"user@gmail.com"}`)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)
}

func TestUpdatingToAlreadyExistingEmail(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 400 trying to update to already existing email"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	userOne := models.User{Name: "username1", Email: "user1@gmail.com", Password: "password"}
	userTwo := models.User{Name: "username2", Email: "user2@gmail.com", Password: "password"}

	userIdOne, err := userOne.Create(userOne)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
		return
	}

	if _, err := userTwo.Create(userTwo); err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
		return
	}

	signInPayload := []byte(`{"email":"user1@gmail.com","password":"password"}`)
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

	path := fmt.Sprintf("/api/v1/auth/user/update/%s", userIdOne)

	payload = []byte(`{"name":"username","email":"user2@gmail.com"}`)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

}

func TestSuccessfulUserDetailUpdate(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 200 on successful user details update"
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

	path := fmt.Sprintf("/api/v1/auth/user/update/%s", userId)

	payload = []byte(`{"name":"username2","email":"user2@gmail.com"}`)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusOK, response.Code)
}
