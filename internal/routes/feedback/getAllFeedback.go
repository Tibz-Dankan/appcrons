package feedback

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getAllFeedback(w http.ResponseWriter, r *http.Request) {
	feedback := models.Feedback{}
	date := services.Date{}

	before := r.URL.Query().Get("before")

	createdAtBefore, err := date.FormatDateString(before)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	allFeedback, count, err := feedback.FindAll(createdAtBefore)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	data := map[string]interface{}{
		"feedback": allFeedback,
		"count":    count,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Fetched user feedback successfully",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetAllFeedbackRoute(router *mux.Router) {
	router.HandleFunc("/get-all", getAllFeedback).Methods("GET")
}
