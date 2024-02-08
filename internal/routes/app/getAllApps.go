package app

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getAllApps(w http.ResponseWriter, r *http.Request) {
	app := models.App{}

	apps, err := app.FindAll()
	if err != nil {
		services.AppError(err.Error(), 400, w)
	}

	data := map[string]interface{}{
		"apps": apps,
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "Apps fetched",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetAllAppsRoute(router *mux.Router) {
	router.HandleFunc("/get-all", getAllApps).Methods("GET")
}
