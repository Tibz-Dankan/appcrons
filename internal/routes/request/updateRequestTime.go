package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func validateUpdateRequestTime(rtId string, startTime, endTime time.Time, savedRequestTimes []models.RequestTime) error {
	// timeLayout -> "2006-Jan-02 15:04:05"

	for _, rt := range savedRequestTimes {
		if rt.ID == rtId {
			continue
		}
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

	if err := validateUpdateRequestTime(requestTime.ID, parsedStartTime, parsedEndTime, savedRequestTimes); err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	updatedRequestTime, err := requestTime.Update()
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	response := map[string]interface{}{
		"status":      "success",
		"message":     "Request Time updated successfully",
		"requestTime": updatedRequestTime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateRequestTimeRoute(router *mux.Router) {
	router.HandleFunc("/update-request-time/{requestTimeId}", updateRequestTime).Methods("PATCH")
}
