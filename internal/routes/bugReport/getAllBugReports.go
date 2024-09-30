package bugreport

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getAllBugReports(w http.ResponseWriter, r *http.Request) {
	bugReport := models.BugReport{}
	date := services.Date{}

	before := r.URL.Query().Get("before")

	createdAtBefore, err := date.FormatDateString(before)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	allBugReports, count, err := bugReport.FindAll(createdAtBefore)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	data := map[string]interface{}{
		"bugReport": allBugReports,
		"count":     count,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Fetched bug reports successfully",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetAllBugReportsRoute(router *mux.Router) {
	router.HandleFunc("/get-all", getAllBugReports).Methods("GET")
}
