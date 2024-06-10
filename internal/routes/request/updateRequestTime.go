package request

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/event"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func updateRequestTime(w http.ResponseWriter, r *http.Request) {
	requestTimeId := mux.Vars(r)["requestTimeId"]
	requestTime := models.RequestTime{ID: requestTimeId}

	err := json.NewDecoder(r.Body).Decode(&requestTime)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if requestTime.ID == "" {
		services.AppError("Please provide requestTimeId", 400, w)
		return
	}

	if requestTime.AppID == "" || requestTime.Start == "" || requestTime.End == "" || requestTime.TimeZone == "" {
		services.AppError("Please fill out all fields!", 400, w)
		return
	}

	savedRequestTime, err := requestTime.FindOne(requestTimeId)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if savedRequestTime.ID == "" {
		services.AppError("Request Time of provided id does not exist", 404, w)
		return
	}

	updatedRequestTime := savedRequestTime
	updatedRequestTime.Start = requestTime.Start
	updatedRequestTime.End = requestTime.End

	err = updatedRequestTime.Update()
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	response := map[string]interface{}{
		"status":      "success",
		"message":     "Request Time updated successfully",
		"requestTime": updatedRequestTime,
	}

	app := models.App{ID: requestTime.AppID}
	event.EB.Publish("updateApp", app)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateRequestTimeRoute(router *mux.Router) {
	router.HandleFunc("/update-request-time/{requestTimeId}", updateRequestTime).Methods("PATCH")
}
