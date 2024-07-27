package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func forgotPassword(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	user, err = user.FindByEMail(user.Email)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if user.ID == "" {
		services.AppError("We couldn't find user with provided email!", 400, w)
		return
	}

	resetToken, err := user.CreatePasswordResetToken()
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	originURL := services.GetRequestOrigin(r)

	resetURL := originURL + "/auth/reset-password/" + resetToken
	log.Println("Password resetURL  ==> ", resetURL)

	// TODO: To send actual email here

	response := map[string]interface{}{
		"status":  "success",
		"message": "Password Reset link sent to email",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ForgotPasswordRoute(router *mux.Router) {
	router.HandleFunc("/forgot-password", forgotPassword).Methods("POST")
}
