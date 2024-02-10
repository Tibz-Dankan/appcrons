package request

import (
	"encoding/json"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func getRequest(w http.ResponseWriter, r *http.Request) {
	request := models.Request{}

	requestId := mux.Vars(r)["requestId"]

	if requestId == "" {
		services.AppError("Please provide requestId", 400, w)
		return
	}

	request, err := request.FindOne(requestId)
	if err != nil {
		services.AppError(err.Error(), 400, w)
	}

	data := map[string]interface{}{
		"request": request,
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "Request fetched",
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetRequestRoute(router *mux.Router) {
	router.HandleFunc("/get/{requestId}", getRequest).Methods("GET")
}
