package app

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/gorilla/mux"
)

func getApp(w http.ResponseWriter, r *http.Request) {
	app := models.App{}
	// get appId from params

	newApp := map[string]interface{}{
		// "id":              appId,
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetAppRoute(router *mux.Router) {
	router.HandleFunc("/get/{appId}", getApp).Methods("GET")
}
