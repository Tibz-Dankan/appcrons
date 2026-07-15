package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getRequestByUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	app := models.App{}
	request := models.Request{}
	date := services.Date{}

	userId := mux.Vars(r)["userId"]
	appId := mux.Vars(r)["appId"]
	before := r.URL.Query().Get("before")

	createdAtBefore, err := date.FormatDateString(before)
	if err != nil {
		services.AppError(err.Error(), 400, w)
	}

	if appId == "" {
		services.AppError("Please provide appId", 400, w)
		return
	}

	// Get user details
	user, err = user.FindOne(userId)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if user.ID == "" {
		services.AppError("User of provided id doesn't exist", 400, w)
		return
	}

	// Get app details
	app, err = app.FindOne(appId)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if app.ID == "" {
		services.AppError("Application of provided id doesn't exist", 400, w)
		return
	}

	// Get requests
	requests, count, err := request.FindByApp(appId, createdAtBefore)
	if err != nil {
		services.AppError(err.Error(), 400, w)
	}

	user.Password = ""
	user.PasswordResetToken = ""
	user.PasswordResetExpiresAt = time.Time{}

	data := map[string]interface{}{
		"user":     user,
		"app":      app,
		"requests": requests,
		"count":    count,
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "Requests fetched",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetRequestsByAppRoute(router *mux.Router) {
	router.HandleFunc("/users/{userId}/apps/{appId}/requests", getRequestByUser).Methods("GET")
}
