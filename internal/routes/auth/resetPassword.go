package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Tibz-Dankan/keep-active/internal/events"
	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func resetPassword(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		services.AppError(err.Error(), 500, w)
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

	clientIP, _ := r.Context().Value(middlewares.ClientIPKey).(string)
	device := r.Header.Get("x-device")

	log.Println("resetPassword clientIP address: ", clientIP)
	log.Println("resetPassword device address: ", device)

	location, err := services.GetUserLocationByIP(user.ID, clientIP)
	if err != nil {
		log.Println("Error resolving location for session:", err)
	}

	session := models.Session{
		UserID:       user.ID,
		AccessToken:  accessToken,
		GeneratedVia: "reset password",
		Device:       device,
		LocationID:   location.ID,
	}
	if _, err := session.Create(session); err != nil {
		log.Println("Error creating session:", err)
	}

	if os.Getenv("GO_ENV") == "testing" || os.Getenv("GO_ENV") == "staging" {
		permission := models.Permissions{}
		if err := permission.Set(user.ID); err != nil {
			log.Println("Error setting permissions:", err)
		}
	} else {
		events.EB.Publish("permissions", user)
	}

	userMap := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}
	tokenData := map[string]interface{}{
		"token": accessToken,
		"user":  userMap,
	}
	response := map[string]interface{}{
		"status":      "success",
		"message":     "Password reset successfully",
		"accessToken": accessToken,
		"user":        userMap,
		"data":        tokenData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ResetPasswordRoute(router *mux.Router) {
	router.HandleFunc("/reset-password/{resetToken}", resetPassword).Methods("PATCH")
}
