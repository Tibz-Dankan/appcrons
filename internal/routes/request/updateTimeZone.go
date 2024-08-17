package request

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func updateTimeZone(w http.ResponseWriter, r *http.Request) {
	requestTime := models.RequestTime{}

	err := json.NewDecoder(r.Body).Decode(&requestTime)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if requestTime.AppID == "" || requestTime.TimeZone == "" {
		services.AppError("Missing appId or timeZone!", 400, w)
		return
	}

	requestTimes, err := requestTime.UpdateTimeZone(requestTime.TimeZone)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	response := map[string]interface{}{
		"status":      "success",
		"message":     "Time zone updated successfully",
		"requestTime": requestTimes,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateTimeZoneRoute(router *mux.Router) {
	router.HandleFunc("/update-timezone", updateTimeZone).Methods("PATCH")
}
