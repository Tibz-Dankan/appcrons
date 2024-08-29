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

func TestMissingPostAppFields(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 400 with postApp fields"
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

	label = "Expects 400 with missing app name"
	payload = []byte(`{"name":"","url":"myapp.onrender.com/active","requestInterval":"5"}`)
	req, _ = http.NewRequest("POST", "/api/v1/apps/post", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

	label = "Expects 400 with missing app url"
	payload = []byte(`{"name":"myApp","url":"","requestInterval":"5"}`)
	req, _ = http.NewRequest("POST", "/api/v1/apps/post", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

	label = "Expects 400 with missing app request interval"
	payload = []byte(`{"name":"myApp","url":"myapp.onrender.com/active","requestInterval":""}`)
	req, _ = http.NewRequest("POST", "/api/v1/apps/post", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)
}

func TestAlreadyExistingApplicationName(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 400 with already existing application name"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	genData := data.NewGenTestData()
	name := genData.RandomUniqueName()
	email := genData.RandomUniqueEmail()
	password := genData.RandomUniquePassword(8)
	appURL := genData.RandomUniqueURL()
	appName := genData.RandomUniqueAppName()

	user := models.User{Name: name, Email: email, Password: password}

	userId, err := user.Create(user)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new user. Got %v\n", err)
		return
	}

	app := models.App{UserID: userId, Name: appName, URL: appURL, RequestInterval: "5"}

	if _, err := app.Create(app); err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
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

	newApp := struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		RequestInterval string `json:"requestInterval"`
	}{
		Name:            appName,
		URL:             genData.RandomUniqueURL(),
		RequestInterval: "5",
	}

	payload, _ = json.Marshal(newApp)
	req, _ = http.NewRequest("POST", "/api/v1/apps/post", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

}

func TestAlreadyExistingApplicationURL(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 400 with already existing application URL"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	genData := data.NewGenTestData()
	name := genData.RandomUniqueName()
	email := genData.RandomUniqueEmail()
	password := genData.RandomUniquePassword(8)
	appURL := genData.RandomUniqueURL()
	appName := genData.RandomUniqueAppName()

	user := models.User{Name: name, Email: email, Password: password}

	userId, err := user.Create(user)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new user. Got %v\n", err)
		return
	}

	app := models.App{UserID: userId, Name: appName, URL: appURL, RequestInterval: "5"}

	if _, err := app.Create(app); err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
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

	newApp := struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		RequestInterval string `json:"requestInterval"`
	}{
		Name:            genData.RandomUniqueAppName(),
		URL:             appURL,
		RequestInterval: "5",
	}

	payload, _ = json.Marshal(newApp)
	req, _ = http.NewRequest("POST", "/api/v1/apps/post", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

}

func TestSuccessfulApplicationPost(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 201 on successful post"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	genData := data.NewGenTestData()
	name := genData.RandomUniqueName()
	email := genData.RandomUniqueEmail()
	password := genData.RandomUniquePassword(8)
	appURL := genData.RandomUniqueURL()
	appName := genData.RandomUniqueAppName()

	user := models.User{Name: name, Email: email, Password: password}

	_, err := user.Create(user)
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

	newApp := struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		RequestInterval string `json:"requestInterval"`
	}{
		Name:            appName,
		URL:             appURL,
		RequestInterval: "5",
	}

	payload, _ = json.Marshal(newApp)
	req, _ = http.NewRequest("POST", "/api/v1/apps/post", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusCreated, response.Code)

}
