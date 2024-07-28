package feedback

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func updateFeedback(w http.ResponseWriter, r *http.Request) {
	feedbackId := mux.Vars(r)["feedbackId"]
	feedback := models.Feedback{ID: feedbackId}

	err := json.NewDecoder(r.Body).Decode(&feedback)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if feedback.ID == "" {
		services.AppError("Please provide feedback id!", 400, w)
		return
	}

	if feedback.Rating == 0 || feedback.Message == "" {
		services.AppError("Missing rating/message!", 400, w)
		return
	}

	savedFeedback, err := feedback.FindOne(feedbackId)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if savedFeedback.ID == "" {
		services.AppError("Feedback of provided id does not exist!", 404, w)
		return
	}

	currentTime := time.Now()
	location := currentTime.Location().String()
	log.Println("time zone location", location)

	appDate := services.Date{TimeZone: location, ISOStringDate: savedFeedback.CreatedAt.String()}

	currentAppTime, _ := appDate.CurrentTime()
	createdAt, _ := appDate.ISOTime()
	timeDifference := currentAppTime.Sub(createdAt).Hours()

	if timeDifference >= 48 {
		services.AppError("Can't edit feedback beyond 48 hours from  its creation!", 403, w)
		return
	}

	updatedFeedbackTime := savedFeedback
	updatedFeedbackTime.Rating = feedback.Rating
	updatedFeedbackTime.Message = feedback.Message

	err = updatedFeedbackTime.Update()
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	response := map[string]interface{}{
		"status":   "success",
		"message":  "Feedback updated successfully",
		"feedback": updatedFeedbackTime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateFeedbackRoute(router *mux.Router) {
	router.HandleFunc("/update/{feedbackId}", updateFeedback).Methods("PATCH")
}
