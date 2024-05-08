package request

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func postRequestTime(w http.ResponseWriter, r *http.Request) {

	requestTime := models.RequestTime{}

	err := json.NewDecoder(r.Body).Decode(&requestTime)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if requestTime.AppID == "" || requestTime.Start == "" || requestTime.End == "" || requestTime.TimeZone == "" {
		services.AppError("Please fill out all fields!", 400, w)
		return
	}
	// TODO: ensure upcoming End time is greater upcoming start time
	currentTime := time.Now()
	timeLayout := currentTime.Format("14:30")
	log.Println("timeLayout::::", timeLayout)

	parsedStartTime, err := time.Parse(timeLayout, requestTime.Start)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	parsedEndTime, err := time.Parse(timeLayout, requestTime.End)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if !parsedStartTime.Before(parsedEndTime) {
		services.AppError("Start Time can't be greater or equal End Time", 400, w)
	}

	// Fetch all request time for the app in the database
	savedRequestTimes, err := requestTime.FindByApp(requestTime.AppID)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	// TODO: add functionality to ensure upcoming time
	// Interval doesn't fall in any existing interval

	for _, rt := range savedRequestTimes {
		parsedSavedRequestTime, err := time.Parse(timeLayout, rt.End)
		if err != nil {
			services.AppError(err.Error(), 400, w)
			return
		}

		if parsedSavedRequestTime.Before(parsedStartTime) {
			services.AppError(errors.New("start time can't be greater than existing end time").Error(), 400, w)
			return
		}
	}

	log.Println("savedRequestTimes:::: ", savedRequestTimes)

	createdRequestTime, err := requestTime.Create(requestTime)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	response := map[string]interface{}{
		"status":      "success",
		"message":     "Request Time Created successfully",
		"requestTime": createdRequestTime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func PostRequestTimeRoute(router *mux.Router) {
	router.HandleFunc("/post-request-time", postRequestTime).Methods("POST")
}
