package request

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getRequestByUser(w http.ResponseWriter, r *http.Request) {
	request := models.Request{}
	date := services.Date{}

	appId := r.URL.Query().Get("appId")
	before := r.URL.Query().Get("before")

	createdAtBefore, err := date.FormatDateString(before)
	if err != nil {
		services.AppError(err.Error(), 400, w)
	}

	if appId == "" {
		services.AppError("Please provide appId", 400, w)
		return
	}

	requests, count, err := request.FindByApp(appId, createdAtBefore)
	if err != nil {
		services.AppError(err.Error(), 400, w)
	}

	data := map[string]interface{}{
		"requests": requests,
		"count":    count,
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "Requests fetched",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetRequestByUserRoute(router *mux.Router) {
	router.HandleFunc("/get-by-app", getRequestByUser).Methods("GET")
}
