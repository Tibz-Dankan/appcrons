package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func updateUserDetails(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userId"]
	user := models.User{ID: userId}

	if user.ID == "" {
		services.AppError("Please provide userId!", 400, w)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	name := user.Name
	email := user.Email

	if user.Name == "" || user.Email == "" {
		services.AppError("Missing user name or email!", 400, w)
		return
	}

	savedUser, err := user.FindOne(user.ID)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if savedUser.ID == "" {
		services.AppError("We couldn't find user of the provided id!", 404, w)
		return
	}

	if savedUser.Email != user.Email {
		existingUser, err := user.FindByEMail(user.Email)
		if err != nil {
			services.AppError(err.Error(), 400, w)
			return
		}
		if existingUser.ID != "" {
			services.AppError("Can't update to already registered email!", 400, w)
			return
		}
	}

	user = savedUser
	user.Email = email
	user.Name = name

	if err := user.Update(); err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	updatedUser := map[string]interface{}{
		"id":    userId,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "User details updated successfully",
		"user":    updatedUser,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateUserDetailsRoute(router *mux.Router) {
	router.HandleFunc("/user/update/{userId}", updateUserDetails).Methods("PATCH")
}
