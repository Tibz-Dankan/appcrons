package app

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func updateApp(w http.ResponseWriter, r *http.Request) {

	appId := mux.Vars(r)["appId"]
	app := models.App{ID: appId}

	if appId == "" {
		services.AppError("Please provide appId", 500, w)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if app.Name == "" || app.URL == "" || app.RequestInterval == "" {
		services.AppError("Please fill out all fields!", 400, w)
		return
	}

	savedApp, err := app.FindAppDetails(appId)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if savedApp.Name != app.Name {
		savedApp, err := app.FindByName(app.Name)
		if err != nil {
			services.AppError(err.Error(), 500, w)
			return
		}
		if savedApp.Name != "" {
			services.AppError("Can't update to already existing app name!", 400, w)
			return
		}

	}

	if savedApp.URL != app.URL {
		savedApp, err := app.FindByURL(app.URL)
		if err != nil {
			services.AppError(err.Error(), 500, w)
			return
		}
		if savedApp.URL != "" {
			services.AppError("Can't update to already existing app url!", 400, w)
			return
		}
	}

	updatedApp := savedApp
	updatedApp.Name = app.Name
	updatedApp.URL = app.URL
	updatedApp.RequestInterval = app.RequestInterval

	app, err = updatedApp.Update()
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Updated successfully",
		"app":     app,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateAppRoute(router *mux.Router) {
	router.HandleFunc("/update/{appId}", updateApp).Methods("PATCH")
}
