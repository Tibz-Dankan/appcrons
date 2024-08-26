package auth_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tibz-Dankan/keep-active/tests/setup"
)

func TestNonExistingForgotPasswordUser(t *testing.T) {
	setup.ClearAllTables()

	var label string
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder

	label = "Expects 400 for non existing user"
	payload = []byte(`{"name":"username","email":"user@gmail.com","password":"password"}`)
	req, _ = http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(payload))
	_ = setup.ExecuteRequest(req)

	payload = []byte(`{"email":"user20@gmail.com"}`)
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

	label = "Expects 400 for non existing user"
	payload = []byte(`{"name":"username","email":"user@gmail.com","password":"password"}`)
	req, _ = http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(payload))
	_ = setup.ExecuteRequest(req)

	payload = []byte(`{"email":"user@gmail.com"}`)
	req, _ = http.NewRequest("POST", "/api/v1/auth/forgot-password", bytes.NewBuffer(payload))
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusOK, response.Code)
}
