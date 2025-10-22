package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userId"]
	user := models.User{ID: userId}

	if user.ID == "" {
		services.AppError("Please provide userId!", 400, w)
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

	userData := map[string]interface{}{
		"id":    userId,
		"name":  savedUser.Name,
		"email": savedUser.Email,
		"role":  savedUser.Role,
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "User fetched successfully",
		"data":    userData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetUserRoute(router *mux.Router) {
	router.HandleFunc("/user/get/{userId}", getUser).Methods("GET")
}
