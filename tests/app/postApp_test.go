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
	var label string = "Expects 400 with postApp fields"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder

	_, _, _, bearerToken := setup.CreateSignInUser()

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
	var label string = "Expects 400 with already existing application name"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder

	userId, _, _, bearerToken := setup.CreateSignInUser()
	genData := data.NewGenTestData()
	appURL := genData.RandomUniqueURL()
	appName := genData.RandomUniqueAppName()

	app := models.App{UserID: userId, Name: appName, URL: appURL, RequestInterval: "5"}

	if _, err := app.Create(app); err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new app. Got %v\n", err)
	}

	appStruct := struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		RequestInterval string `json:"requestInterval"`
	}{
		Name:            appName,
		URL:             genData.RandomUniqueURL(),
		RequestInterval: "5",
	}

	payload, _ = json.Marshal(appStruct)
	req, _ = http.NewRequest("POST", "/api/v1/apps/post", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

}

func TestAlreadyExistingApplicationURL(t *testing.T) {
	var label string = "Expects 400 with already existing application URL"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder

	userId, _, _, bearerToken := setup.CreateSignInUser()

	genData := data.NewGenTestData()
	appURL := genData.RandomUniqueURL()
	appName := genData.RandomUniqueAppName()

	app := models.App{UserID: userId, Name: appName, URL: appURL, RequestInterval: "5"}

	if _, err := app.Create(app); err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new app. Got %v\n", err)
	}

	appStruct := struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		RequestInterval string `json:"requestInterval"`
	}{
		Name:            genData.RandomUniqueAppName(),
		URL:             appURL,
		RequestInterval: "5",
	}

	payload, _ = json.Marshal(appStruct)
	req, _ = http.NewRequest("POST", "/api/v1/apps/post", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

}

func TestSuccessfulApplicationPost(t *testing.T) {
	var label string = "Expects 201 on successful post"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder

	_, _, _, bearerToken := setup.CreateSignInUser()

	genData := data.NewGenTestData()
	appURL := genData.RandomUniqueURL()
	appName := genData.RandomUniqueAppName()

	appStruct := struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		RequestInterval string `json:"requestInterval"`
	}{
		Name:            appName,
		URL:             appURL,
		RequestInterval: "5",
	}

	payload, _ = json.Marshal(appStruct)
	req, _ = http.NewRequest("POST", "/api/v1/apps/post", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusCreated, response.Code)

}
