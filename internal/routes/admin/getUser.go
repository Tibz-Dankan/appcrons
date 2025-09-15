package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	userId := mux.Vars(r)["userId"]

	if userId == "" {
		services.AppError("Please provide user id", 400, w)
		return
	}

	user, err := user.FindOne(userId)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if user.ID == "" {
		services.AppError("User of the provided id doesn't exist", 404, w)
		return
	}

	user.Password = ""
	user.PasswordResetToken = ""
	user.PasswordResetExpiresAt = time.Time{}

	response := map[string]interface{}{
		"status":  "success",
		"message": "User fetched successfully",
		"data":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetUserRoute(router *mux.Router) {
	router.HandleFunc("/users/{userId}", getUser).Methods("GET")
}