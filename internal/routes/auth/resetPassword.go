package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/events"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func resetPassword(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	newPassword := user.Password
	token := mux.Vars(r)["resetToken"]

	if user.Password == "" {
		services.AppError("Please provide your new password!", 400, w)
		return
	}

	user, err = user.FindByPasswordResetToken(token)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if user.ID == "" {
		services.AppError("Invalid or expired reset token!", 400, w)
		return
	}

	err = user.ResetPassword(newPassword)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	accessToken, err := services.SignJWTToken(user.ID)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	userMap := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}
	response := map[string]interface{}{
		"status":      "success",
		"message":     "Password reset successfully",
		"accessToken": accessToken,
		"user":        userMap,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	events.EB.Publish("permissions", user)
}

func ResetPasswordRoute(router *mux.Router) {
	router.HandleFunc("/reset-password/{resetToken}", resetPassword).Methods("PATCH")
}
