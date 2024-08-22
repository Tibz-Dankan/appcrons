package request

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func deleteRequestTime(w http.ResponseWriter, r *http.Request) {
	requestTimeId := mux.Vars(r)["requestTimeId"]
	requestTime := models.RequestTime{ID: requestTimeId}

	if requestTime.ID == "" {
		services.AppError("Please provide requestTimeId", 400, w)
		return
	}

	savedRequestTime, err := requestTime.FindOne(requestTime.ID)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if savedRequestTime.ID == "" {
		services.AppError("Request Time of provided id doesn't exist", 404, w)
		return
	}

	err = requestTime.Delete(requestTime.ID)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Request Time deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteRequestTimeRoute(router *mux.Router) {
	router.HandleFunc("/delete-request-time/{requestTimeId}", deleteRequestTime).Methods("DELETE")
}
