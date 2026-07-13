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

func signUp(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		services.AppError("Please fill out all fields!", 400, w)
		return
	}

	savedUser, err := user.FindByEmail(user.Email)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if savedUser.ID != "" {
		services.AppError("Email already registered!", 400, w)
		return
	}

	err = user.SetRole("user")
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	userId, err := user.Create(user)

	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	accessToken, err := services.SignJWTToken(userId)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	user.ID = userId

	clientIP, _ := r.Context().Value(middlewares.ClientIPKey).(string)
	device := r.Header.Get("User-Agent")

	location, err := services.GetUserLocationByIP(user.ID, clientIP)
	if err != nil {
		log.Println("Error resolving location for session:", err)
	}

	session := models.Session{
		UserID:       user.ID,
		AccessToken:  accessToken,
		GeneratedVia: "sign up",
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

	newUser := map[string]interface{}{
		"id":    userId,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}
	tokenData := map[string]interface{}{
		"token": accessToken,
		"user":  newUser,
	}
	response := map[string]interface{}{
		"status":      "success",
		"message":     "Signup successfully",
		"accessToken": accessToken,
		"user":        newUser,
		"data":        tokenData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func SignUpRoute(router *mux.Router) {
	router.HandleFunc("/signup", signUp).Methods("POST")
}
