package app

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/event"
	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func PostAdd(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		services.AppError("UserID not found in context", 500, w)
		return
	}

	app := models.App{UserID: userID}

	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if app.Name == "" || app.URL == "" || app.RequestInterval == "" {
		services.AppError("Please fill out all fields!", 400, w)
		return
	}

	savedApp, err := app.FindByName(app.Name)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if savedApp.ID != "" {
		services.AppError("App name already taken!", 400, w)
		return
	}

	savedApp, err = app.FindByURL(app.URL)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if savedApp.URL != "" {
		services.AppError("App URL already exists!", 400, w)
		return
	}

	createdApp, err := app.Create(app)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	event.EB.Publish("createApp", app)

	response := map[string]interface{}{
		"status":  "success",
		"message": "Created successfully",
		"app":     createdApp,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func PostAppRoute(router *mux.Router) {
	router.HandleFunc("/post", PostAdd).Methods("POST")
}
