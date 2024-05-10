package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

var timeLayout string = "2006-Jan-02 15:04:05"

func timeValue(timeArg string) string {
	currentTime := time.Now()

	var currentTimeDay string
	day := currentTime.Day()
	if day < 10 {
		currentTimeDay = "0" + fmt.Sprint(day)
	} else {
		currentTimeDay = fmt.Sprint(day)
	}
	currentTimeYear := fmt.Sprint(currentTime.Year())
	currentTimeMonth := fmt.Sprint(currentTime.Month())

	date := currentTimeYear + "-" + currentTimeMonth + "-" + currentTimeDay
	dateValueStr := date + " " + timeArg

	return dateValueStr
}

func validateRequestTime(startTime, endTime time.Time, savedRequestTimes []models.RequestTime) error {
	// timeLayout -> "2006-Jan-02 15:04:05"

	for _, rt := range savedRequestTimes {
		savedEndTime, err := time.Parse(timeLayout, timeValue(rt.End))
		if err != nil {
			return err
		}

		savedStartTime, err := time.Parse(timeLayout, timeValue(rt.Start))
		if err != nil {
			return err
		}

		if savedStartTime.Equal(startTime) || savedEndTime.Equal(endTime) {
			return errors.New("time intervals can't be equal to existing ones")
		}

		// Check for upcoming interval top
		if savedStartTime.Equal(endTime) {
			return errors.New("end time can't be equal to any existing interval")
		}

		if savedStartTime.After(endTime) {
			if !savedStartTime.After(startTime) {
				return errors.New("start time can't be greater than existing start time")
			}
		}

		// Check for upcoming interval down
		if savedEndTime.Equal(startTime) {
			return errors.New("start time can't be equal to any existing interval")
		}

		if savedEndTime.Before(startTime) {
			if !savedEndTime.Before(endTime) {
				return errors.New("end time can't be greater than existing end time")
			}
		}
	}
	return nil
}

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

	parsedStartTime, err := time.Parse(timeLayout, timeValue(requestTime.Start))
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	parsedEndTime, err := time.Parse(timeLayout, timeValue(requestTime.End))
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

	if err := validateRequestTime(parsedStartTime, parsedEndTime, savedRequestTimes); err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

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
