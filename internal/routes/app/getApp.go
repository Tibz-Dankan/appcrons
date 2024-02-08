package app

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getApp(w http.ResponseWriter, r *http.Request) {
	app := models.App{}

	appId := mux.Vars(r)["appId"]

	app, err := app.FindOne(appId)
	if err != nil {
		services.AppError(err.Error(), 500, w)
	}

	data := map[string]interface{}{
		"app": app,
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "App fetched",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetAppRoute(router *mux.Router) {
	router.HandleFunc("/get/{appId}", getApp).Methods("GET")
}
