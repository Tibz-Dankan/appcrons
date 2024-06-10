package request

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/event"
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

	currentTime := time.Now()
	location := currentTime.Location().String()

	appDate := services.Date{TimeZone: location, ISOStringDate: currentTime.String(), HourMinSec: requestTime.Start}

	parsedStartTime, err := appDate.HourMinSecTime()
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}
	log.Println("parsedStartTime:::: ", parsedStartTime)

	// TODO: to validate the requestTime being submitted

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

	app := models.App{ID: requestTime.AppID}
	event.EB.Publish("updateApp", app)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func PostRequestTimeRoute(router *mux.Router) {
	router.HandleFunc("/post-request-time", postRequestTime).Methods("POST")
}
