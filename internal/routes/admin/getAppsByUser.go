package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getAppByUser(w http.ResponseWriter, r *http.Request) {
	app := models.App{}
	user := models.User{}

	userId := mux.Vars(r)["userId"]
	limitParam := r.URL.Query().Get("limit")
	cursorParam := r.URL.Query().Get("cursor")

	limit, err := services.ValidateQueryLimit(limitParam)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if cursorParam == "" {
		cursorParam = ""
	}

	if userId == "" {
		services.AppError("Please provide userId", 400, w)
		return
	}

	user, err = user.FindOne(userId)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if user.ID == "" {
		services.AppError("User of provided id doesn't exist", 400, w)
		return
	}

	apps, err := app.FindByUserPaginated(userId, limit, cursorParam)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	user.Password = ""
	user.PasswordResetToken = ""
	user.PasswordResetExpiresAt = time.Time{}

	data := map[string]interface{}{
		"user": user,
		"apps": apps,
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "Apps fetched",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetAppsByUserRoute(router *mux.Router) {
	router.HandleFunc("/users/{userId}/apps", getAppByUser).Methods("GET")
}
