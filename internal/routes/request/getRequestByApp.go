package request

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getRequestByUser(w http.ResponseWriter, r *http.Request) {
	request := models.Request{}

	appId := r.URL.Query().Get("appId")
	before := r.URL.Query().Get("before")

	var createdAtBefore time.Time

	if before == "" {
		createdAtBefore = time.Now()
		log.Println("createdAtBefore: ", createdAtBefore)
	}

	if before != "" {
		log.Println("before: ", before)

		log.Println("beforeWithReplaced spaces: ", services.ReplaceSpaces(before))

		// Check if before contains Z (UTC format)
		isUTC := services.ContainsZ(services.ReplaceSpaces(before))

		if isUTC {
			date := services.Date{ISOStringDate: services.ReplaceSpaces(before)}

			createdAt, err := date.UTC()
			if err != nil {
				services.AppError("Something went wrong, please try again later!", 500, w)
				return
			}
			createdAtBefore = createdAt
			log.Println("before: ", before)
			log.Println("createdAtBefore: ", createdAtBefore)
		}

		if !isUTC {
			date := services.Date{ISOStringDate: services.ReplaceSpaces(before)}

			createdAt, err := date.RFC3339Nano()
			if err != nil {
				services.AppError("Something went wrong, please try again later!", 500, w)
				return
			}
			createdAtBefore = createdAt
			log.Println("before: ", before)
			log.Println("createdAtBefore: ", createdAtBefore)
		}
	}

	if appId == "" {
		services.AppError("Please provide appId", 400, w)
		return
	}

	requests, count, err := request.FindByApp(appId, createdAtBefore)
	if err != nil {
		services.AppError(err.Error(), 400, w)
	}

	data := map[string]interface{}{
		"requests": requests,
		"count":    count,
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "Requests fetched",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetRequestByUserRoute(router *mux.Router) {
	router.HandleFunc("/get-by-app", getRequestByUser).Methods("GET")
}
