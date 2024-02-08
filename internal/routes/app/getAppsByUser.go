package app

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/gorilla/mux"
)

func getAppByUser(w http.ResponseWriter, r *http.Request) {
	app := models.App{}
	// get userId from query params

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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetAppByUserByUser(router *mux.Router) {
	router.HandleFunc("/get-by-user", getAppByUser).Methods("GET")
}
