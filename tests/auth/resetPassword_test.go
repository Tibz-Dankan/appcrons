package tests

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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
	"github.com/google/uuid"
)

func TestMissingResetPasswordToken(t *testing.T) {
	setup.ClearAllTables()

	var label string
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder

	label = "Expects 400 with missing reset password token"

	payload = []byte(`{"password":"password"}`)
	req, _ = http.NewRequest("PATCH", "/api/v1/auth/reset-password/''", bytes.NewBuffer(payload))
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)
}

func TestMissingNewResetPassword(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 400 with expired password reset token"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder

	genData := data.NewGenTestData()
	name := genData.RandomUniqueName()
	email := genData.RandomUniqueEmail()
	password := genData.RandomUniquePassword(8)

	user := models.User{Name: name, Email: email, Password: password}

	userId, err := user.Create(user)
	if err != nil {
		log.Println(err)
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
		return
	}
	user.ID = userId
	time.Sleep(500 * time.Millisecond)

	resetToken, err := user.CreatePasswordResetToken()
	if err != nil {
		log.Println(err)
		return
	}

	payload = []byte(`{"password":""}`)
	path := fmt.Sprintf("/api/v1/auth/reset-password/%s", resetToken)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)
}

func TestExpiredPasswordResetToken(t *testing.T) {
	setup.ClearAllTables()

	var label string
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder

	db := setup.DB

	genData := data.NewGenTestData()
	name := genData.RandomUniqueName()
	email := genData.RandomUniqueEmail()
	password := genData.RandomUniquePassword(8)

	user := models.User{Name: name, Email: email, Password: password}

	userId, err := user.Create(user)
	if err != nil {
		log.Println(err)
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
		return
	}
	user.ID = userId
	time.Sleep(500 * time.Millisecond)

	resetToken := uuid.NewString()
	hashedToken := sha256.New()
	hashedToken.Write([]byte(resetToken))
	hashedTokenByteSlice := hashedToken.Sum(nil)
	hashedTokenString := hex.EncodeToString(hashedTokenByteSlice)
	passwordResetExpiresAt := time.Now().Add(-20 * time.Minute)

	stmt := fmt.Sprintf("UPDATE users SET \"passwordResetToken\" = %s \"passwordResetExpiresAt\" = %v  WHERE id = %s", hashedTokenString, passwordResetExpiresAt, userId)
	db.Exec(stmt)

	label = "Expects 400 with expired password reset token"
	payload = []byte(`{"password":"newPassword"}`)
	path := fmt.Sprintf("/api/v1/auth/reset-password/%s", resetToken)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusBadRequest, response.Code)
}

func TestSuccessfulPasswordReset(t *testing.T) {
	setup.ClearAllTables()

	var label string = "Expects 200 on successful password reset"
	var payload []byte
	var req *http.Request
	var response *httptest.ResponseRecorder
	var body map[string]interface{}

	user := models.User{Name: "username", Email: "user@gmail.com", Password: "password"}

	userId, err := user.Create(user)
	if err != nil {
		log.Println(err)
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects accessToken. Got %v\n", err)
		return
	}
	user.ID = userId

	resetToken, err := user.CreatePasswordResetToken()
	if err != nil {
		fmt.Printf("=== FAIL: %s\n", label)
		t.Errorf("Expects resetToken. Got %v\n", err)
	}

	payload = []byte(`{"password":"newPassword"}`)
	path := fmt.Sprintf("/api/v1/auth/reset-password/%s", resetToken)
	req, _ = http.NewRequest("PATCH", path, bytes.NewBuffer(payload))
	response = setup.ExecuteRequest(req)
	setup.CheckResponseCode(t, label, http.StatusOK, response.Code)

	json.Unmarshal(response.Body.Bytes(), &body)

	label = "Expects accessToken on successful password reset"
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
