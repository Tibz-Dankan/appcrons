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

func signIn(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	password := user.Password

	if user.Email == "" || user.Password == "" {
		services.AppError("Missing email/password!", 400, w)
		return
	}

	user, err = user.FindByEmail(user.Email)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if user.ID == "" {
		services.AppError("Invalid email/password!", 400, w)
		return
	}

	passwordMatches, err := user.PasswordMatches(password)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if !passwordMatches {
		services.AppError("Invalid email/password!", 400, w)
		return
	}

	accessToken, err := services.SignJWTToken(user.ID)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	clientIP, _ := r.Context().Value(middlewares.ClientIPKey).(string)
	device := r.Header.Get("User-Agent")

	location, err := services.GetUserLocationByIP(user.ID, clientIP)
	if err != nil {
		log.Println("Error resolving location for session:", err)
	}

	session := models.Session{
		UserID:       user.ID,
		AccessToken:  accessToken,
		GeneratedVia: "sign in",
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
		"message":     "Sign in successfully",
		"accessToken": accessToken,
		"user":        userMap,
		"data":        tokenData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SignInRoute(router *mux.Router) {
	router.HandleFunc("/signin", signIn).Methods("POST")
}
