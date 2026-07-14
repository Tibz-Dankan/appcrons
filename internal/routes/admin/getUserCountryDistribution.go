package admin

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getUserCountryDistribution(w http.ResponseWriter, r *http.Request) {
	location := models.Location{}
	user := models.User{}

	distribution, err := location.FindUserCountryDistribution()
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	totalUserCount, err := user.FindCount()
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	var knownUserCount int64
	unknownIndex := -1
	for i, d := range distribution {
		knownUserCount += d.UserCount
		if d.CountryCode == "" {
			unknownIndex = i
		}
	}

	if usersWithoutLocation := totalUserCount - knownUserCount; usersWithoutLocation > 0 {
		if unknownIndex >= 0 {
			distribution[unknownIndex].UserCount += usersWithoutLocation
		} else {
			distribution = append(distribution, models.CountryDistribution{
				Country: "Unknown", CountryCode: "", UserCount: usersWithoutLocation,
			})
		}
	}

	sort.Slice(distribution, func(i, j int) bool {
		return distribution[i].UserCount > distribution[j].UserCount
	})

	response := map[string]interface{}{
		"status":  "success",
		"message": "User country distribution fetched successfully",
		"data":    distribution,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetUserCountryDistributionRoute(router *mux.Router) {
	router.HandleFunc("/users/countries", getUserCountryDistribution).Methods("GET")
}
