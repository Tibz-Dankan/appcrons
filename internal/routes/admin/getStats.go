package admin

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getStats(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	app := models.App{}
	request := models.Request{}
	location := models.Location{}


	userCount, err := user.FindCount()
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	appCount, err := app.FindCount()
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	requestCount, err := request.FindTotalCount()
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	countryDistribution, err := location.FindUserCountryDistribution()
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	countryCount := 0
	for _, d := range countryDistribution {
		if d.CountryCode != "" {
			countryCount++
		}
	}

	data := map[string]interface{}{
		"userCount":    userCount,
		"appCount":     appCount,
		"requestCount": requestCount,
		"countryCount": countryCount,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Stats fetched successfully",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetStatsRoute(router *mux.Router) {
	router.HandleFunc("/stats", getStats).Methods("GET")
}