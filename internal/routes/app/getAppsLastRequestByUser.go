package app

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getAppsLastRequestByUser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")

	app := models.App{UserID: userId}

	if userId == "" {
		services.AppError("Please provide user id", 400, w)
		return
	}

	data := map[string]interface{}{
		"apps": []services.AppRequestProgress{},
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "User applications fetched successfully",
		"data":    data,
	}

	// Get apps from the memory and send them to the client
	if apps, found := services.UserAppMem.Get(userId); found {
		data["apps"] = apps

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get apps from the database
	apps, err := app.FindByUser(userId)
	if err != nil {
		services.AppError(err.Error(), 500, w)
	}

	appRequestProgressList := []services.AppRequestProgress{}

	for _, app := range apps {
		appRequestProgress := services.AppRequestProgress{App: app, InProgress: false}
		appRequestProgressList = append(appRequestProgressList, appRequestProgress)
	}

	// data["apps"] = apps
	data["apps"] = appRequestProgressList

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetAppsLastRequestByUserRoute(router *mux.Router) {
	router.HandleFunc("/get-apps-last-request-by-user", getAppsLastRequestByUser).Methods("GET")
}
