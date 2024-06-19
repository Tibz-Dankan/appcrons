package app

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func deleteApp(w http.ResponseWriter, r *http.Request) {
	appId := mux.Vars(r)["appId"]
	app := models.App{ID: appId}

	if app.ID == "" {
		services.AppError("Please provide requestTimeId", 400, w)
		return
	}

	savedApp, err := app.FindOne(app.ID)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if savedApp.ID == "" {
		services.AppError("App of provided id doesn't exist", 404, w)
		return
	}

	err = app.Delete(app.ID)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": savedApp.Name + " deleted successfully",
	}

	// event.EB.Publish("delete", app)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteAppRoute(router *mux.Router) {
	router.HandleFunc("/delete/{appId}", deleteApp).Methods("DELETE")
}
