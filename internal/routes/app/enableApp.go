package app

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/event"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func enableApp(w http.ResponseWriter, r *http.Request) {

	appId := mux.Vars(r)["appId"]
	app := models.App{ID: appId}

	if appId == "" {
		services.AppError("Please provide appId", 500, w)
	}

	savedApp, err := app.FindAppDetails(appId)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if !savedApp.IsDisabled {
		services.AppError("app is already enabled", 400, w)
		return
	}

	app = savedApp
	app.IsDisabled = false

	app, err = app.Update()
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	event.EB.Publish("updateApp", app)

	response := map[string]interface{}{
		"status":  "success",
		"message": "App is enabled successfully",
		"app":     app,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func EnableAppRoute(router *mux.Router) {
	router.HandleFunc("/enable/{appId}", enableApp).Methods("PATCH")
}
