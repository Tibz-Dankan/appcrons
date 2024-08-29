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

func TestMissingUpdateAppFields(t *testing.T) {
	var label string = "Expects 400 with update app fields"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}
	var userId string
	var bearerToken string
	var email string
	var password string

	userId, email, password, _ = setup.CreateSignInUser()

	genData := data.NewGenTestData()
	appURL := genData.RandomUniqueURL()
	appName := genData.RandomUniqueAppName()

	app := models.App{UserID: userId, Name: appName, URL: appURL, RequestInterval: "5"}

	newApp, err := app.Create(app)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new app. Got %v\n", err)
	}
	app.ID = newApp.ID

	_, bearerToken = setup.SignInUser(email, password)

	appStruct1 := struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		RequestInterval string `json:"requestInterval"`
	}{
		Name:            "",
		URL:             genData.RandomUniqueURL(),
		RequestInterval: "5",
	}
	payload, _ = json.Marshal(appStruct1)

	path := fmt.Sprintf("/api/v1/apps/update/%s", app.ID)

	label = "Expects 400 with missing app name"
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

	json.Unmarshal(response.Body.Bytes(), &body)

	appStruct2 := struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		RequestInterval string `json:"requestInterval"`
	}{
		Name:            appName,
		URL:             "",
		RequestInterval: "5",
	}
	payload, _ = json.Marshal(appStruct2)

	label = "Expects 400 with missing update app URL"
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

	appStruct3 := struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		RequestInterval string `json:"requestInterval"`
	}{
		Name:            appName,
		URL:             appURL,
		RequestInterval: "",
	}
	payload, _ = json.Marshal(appStruct3)

	label = "Expects 400 with missing update app request interval"
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

}

func TestAlreadyExistingUpdateApplicationName(t *testing.T) {
	var label string = "Expects 400 with already existing update application name"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var userId string
	var bearerToken string
	var email string
	var password string

	userId, email, password, _ = setup.CreateSignInUser()
	genData := data.NewGenTestData()

	app1 := models.App{
		UserID: userId, Name: genData.RandomUniqueAppName(),
		URL:             genData.RandomUniqueURL(),
		RequestInterval: "5",
	}
	app2 := models.App{
		UserID: userId, Name: genData.RandomUniqueAppName(),
		URL:             genData.RandomUniqueURL(),
		RequestInterval: "5",
	}

	newApp1, err := app1.Create(app1)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new app. Got %v\n", err)
	}
	app1.ID = newApp1.ID

	if _, err := app2.Create(app2); err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new app. Got %v\n", err)
	}

	_, bearerToken = setup.SignInUser(email, password)

	appStruct := struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		RequestInterval string `json:"requestInterval"`
	}{
		Name:            app2.Name,
		URL:             genData.RandomUniqueURL(),
		RequestInterval: "5",
	}
	payload, _ = json.Marshal(appStruct)

	path := fmt.Sprintf("/api/v1/apps/update/%s", app1.ID)

	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

}

func TestAlreadyExistingUpdateApplicationURL(t *testing.T) {
	var label string = "Expects 400 with already existing update application URL"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var userId string
	var bearerToken string
	var email string
	var password string

	userId, email, password, _ = setup.CreateSignInUser()
	genData := data.NewGenTestData()

	app1 := models.App{
		UserID: userId, Name: genData.RandomUniqueAppName(),
		URL:             genData.RandomUniqueURL(),
		RequestInterval: "5",
	}
	app2 := models.App{
		UserID: userId, Name: genData.RandomUniqueAppName(),
		URL:             genData.RandomUniqueURL(),
		RequestInterval: "5",
	}

	newApp1, err := app1.Create(app1)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new app. Got %v\n", err)
	}
	app1.ID = newApp1.ID

	if _, err := app2.Create(app2); err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects new app2. Got %v\n", err)
	}

	_, bearerToken = setup.SignInUser(email, password)

	appStruct := struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		RequestInterval string `json:"requestInterval"`
	}{
		Name:            genData.RandomUniqueAppName(),
		URL:             app2.URL,
		RequestInterval: "5",
	}
	payload, _ = json.Marshal(appStruct)

	path := fmt.Sprintf("/api/v1/apps/update/%s", app1.ID)

	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)

}

func TestSuccessfulUpdateApp(t *testing.T) {
	var label string = "Expects 400 with already existing update application URL"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var userId string
	var bearerToken string
	var email string
	var password string

	userId, email, password, _ = setup.CreateSignInUser()

	genData := data.NewGenTestData()
	appURL := genData.RandomUniqueURL()
	appName := genData.RandomUniqueAppName()

	app := models.App{UserID: userId, Name: appName, URL: appURL, RequestInterval: "5"}

	newApp, err := app.Create(app)
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
	}
	app.ID = newApp.ID

	_, bearerToken = setup.SignInUser(email, password)

	appVariant1 := struct {
		Name            string `json:"name"`
		URL             string `json:"url"`
		RequestInterval string `json:"requestInterval"`
	}{
		Name:            genData.RandomUniqueAppName(),
		URL:             genData.RandomUniqueURL(),
		RequestInterval: "5",
	}
	payload, _ = json.Marshal(appVariant1)

	path := fmt.Sprintf("/api/v1/apps/update/%s", app.ID)

	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", bearerToken)
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusOK, response.Code)

}
