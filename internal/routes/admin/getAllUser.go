package admin

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

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

	usersWithAppCount, err := user.FindAllAndIncludeAppCount(limit, cursorParam)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	var prevCursor string
	if len(usersWithAppCount) > 0 {
		prevCursor = usersWithAppCount[len(usersWithAppCount)-1].ID
	}

	pagination := map[string]interface{}{
		"limit":      limit,
		"prevCursor": prevCursor,
	}

	response := map[string]interface{}{
		"status":     "success",
		"message":    "Users fetched successfully",
		"data":       usersWithAppCount,
		"pagination": pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetAllUsersRoute(router *mux.Router) {
	router.HandleFunc("/users", getAllUsers).Methods("GET")
}
