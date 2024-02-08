package app

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getAppByUser(w http.ResponseWriter, r *http.Request) {
	app := models.App{}

	userId := r.URL.Query().Get("userId")

	if userId == "" {
		services.AppError("Please provide userId", 400, w)
		return
	}

	apps, err := app.FindByUser(userId)
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

func GetAppByUserByUser(router *mux.Router) {
	router.HandleFunc("/get-by-user", getAppByUser).Methods("GET")
}
