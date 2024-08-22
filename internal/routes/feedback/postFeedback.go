package feedback

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/events"
	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func postFeedback(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		services.AppError("UserID not found in context!", 500, w)
		return
	}

	feedback := models.Feedback{UserID: userId}

	err := json.NewDecoder(r.Body).Decode(&feedback)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if feedback.Rating == 0 || feedback.Message == "" {
		services.AppError("Missing rating/message!", 400, w)
		return
	}

	newFeedback, err := feedback.Create(feedback)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	response := map[string]interface{}{
		"status":   "success",
		"message":  "Thank very much for your feedback",
		"feedback": newFeedback,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	user := models.User{ID: userId}
	events.EB.Publish("permissions", user)
}

func PostFeedbackRoute(router *mux.Router) {
	router.HandleFunc("/post", postFeedback).Methods("POST")
}
