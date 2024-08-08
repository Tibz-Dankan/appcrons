package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/event"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func signUp(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		services.AppError("Please fill out all fields!", 400, w)
		return
	}

	savedUser, err := user.FindByEMail(user.Email)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if savedUser.ID != "" {
		services.AppError("Email already registered!", 400, w)
		return
	}

	// TODO: To implement more scalable approach to set user roles
	err = user.SetRole("client")
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

	newUser := map[string]interface{}{
		"id":    userId,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}
	response := map[string]interface{}{
		"status":      "success",
		"message":     "Signup successfully",
		"accessToken": accessToken,
		"user":        newUser,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	user.ID = userId
	event.EB.Publish("permissions", user)
}

func SignUpRoute(router *mux.Router) {
	router.HandleFunc("/signup", signUp).Methods("POST")
}
