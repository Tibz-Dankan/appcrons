package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

type Passwords struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userId"]
	user := models.User{ID: userId}
	passwords := Passwords{}

	err := json.NewDecoder(r.Body).Decode(&passwords)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if passwords.CurrentPassword == "" || passwords.NewPassword == "" {
		services.AppError("Missing current/new password!", 400, w)
		return
	}

	if passwords.CurrentPassword == passwords.NewPassword {
		services.AppError("New password is same as current password", 400, w)
		return
	}

	user, err = user.FindOne(user.ID)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if user.ID == "" {
		services.AppError("We couldn't find user of provided id!", 400, w)
		return
	}

	currentPasswordMatchesSavedOne, err := user.PasswordMatches(passwords.CurrentPassword)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if !currentPasswordMatchesSavedOne {
		services.AppError("Invalid current password!", 400, w)
		return
	}

	hashedPassword, err := user.HashPassword(passwords.NewPassword)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	user.Password = hashedPassword

	if err := user.Update(); err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Password changed successfully!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ChangePasswordRoute(router *mux.Router) {
	router.HandleFunc("/user/update-password/{userId}", changePassword).Methods("PATCH")
}
