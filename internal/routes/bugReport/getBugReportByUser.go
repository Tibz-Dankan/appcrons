package bugreport

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getBugReportByUser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	bugReport := models.BugReport{UserID: userId}
	date := services.Date{}

	before := r.URL.Query().Get("before")

	createdAtBefore, err := date.FormatDateString(before)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if userId == "" {
		services.AppError("Please provide appId", 400, w)
		return
	}

	userBugReports, count, err := bugReport.FindByUser(userId, createdAtBefore)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	data := map[string]interface{}{
		"bugReport": userBugReports,
		"count":     count,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Fetched user bug reports successfully",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetBugReportByUserRoute(router *mux.Router) {
	router.HandleFunc("/get-by-user", getBugReportByUser).Methods("GET")
}
