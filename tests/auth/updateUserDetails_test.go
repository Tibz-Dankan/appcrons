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

func TestMissingUpdateUserDetailFields(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 400 with missing name/email"
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
	time.Sleep(500 * time.Millisecond)

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

	genData := data.NewGenTestData()

	userOne := models.User{
		Name:     genData.RandomUniqueAppName(),
		Email:    genData.RandomUniqueEmail(),
		Password: genData.RandomUniquePassword(8),
	}
	userTwo := models.User{
		Name:     genData.RandomUniqueAppName(),
		Email:    genData.RandomUniqueEmail(),
		Password: genData.RandomUniquePassword(8),
	}

	userIdOne, err := userOne.Create(userOne)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new user 1. Got %v\n", err)
		return
	}
	userOne.ID = userIdOne

	if _, err := userTwo.Create(userTwo); err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new user 2. Got %v\n", err)
		return
	}
	time.Sleep(500 * time.Millisecond)

	signInPayload, _ := json.Marshal(userOne)
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

	userOne.Email = userTwo.Email
	payload, _ = json.Marshal(userOne)
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
	time.Sleep(500 * time.Millisecond)

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

	path := fmt.Sprintf("/api/v1/auth/user/update/%s", userId)

	user.Email = genData.RandomUniqueEmail()
	user.Name = genData.RandomUniqueName()

	payload, _ = json.Marshal(user)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusOK, response.Code)
}
