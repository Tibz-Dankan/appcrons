package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Tibz-Dankan/keep-active/internal/events"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func signUpAdmin(w http.ResponseWriter, r *http.Request) {
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

	adminEmail := os.Getenv("ADMIN_EMAIL_APPCRONS")
	if adminEmail == "" {
		services.AppError("Invalid admin email", 400, w)
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

	err = user.SetRole("sys_admin")
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
	response := map[string]interface{}{
		"status":      "success",
		"message":     "Signup successfully",
		"accessToken": accessToken,
		"user":        newUser,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func SignUpAdminRoute(router *mux.Router) {
	router.HandleFunc("/signup-admin", signUpAdmin).Methods("POST")
}
