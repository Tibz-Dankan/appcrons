package app

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func PostAdd(w http.ResponseWriter, r *http.Request) {
	app := models.App{}

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

	if savedApp.URL != "" {
		services.AppError("App URL already exists!", 400, w)
		return
	}

	appId, err := app.Create(app)

	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	newApp := map[string]interface{}{
		"id":              appId,
		"name":            app.Name,
		"requestInterval": app.RequestInterval,
		"updatedAt":       app.UpdatedAt,
		"createdAt":       app.CreatedAt,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Created successfully",
		"app":     newApp,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func PostAppRoute(router *mux.Router) {
	router.HandleFunc("/api/v1/app/post", PostAdd).Methods("POST")
}
