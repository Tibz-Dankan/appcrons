package app

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func searchApps(w http.ResponseWriter, r *http.Request) {
	app := models.App{}

	userId := r.URL.Query().Get("userId")
	query := r.URL.Query().Get("query")

	if userId == "" || query == "" {
		services.AppError("Missing query or user id", 400, w)
		return
	}

	apps, err := app.Search(query, userId)
	if err != nil {
		services.AppError(err.Error(), 400, w)
	}

	data := map[string]interface{}{
		"apps": apps,
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "Search results for" + "'" + query + "'",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SearchAppsRoute(router *mux.Router) {
	router.HandleFunc("/search", searchApps).Methods("GET")
}
