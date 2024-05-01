package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func getActive(w http.ResponseWriter, r *http.Request) {

	response := map[string]interface{}{
		"status":  "success",
		"message": "Active",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetActiveRoute(router *mux.Router) {
	router.HandleFunc("/get/active", getActive).Methods("GET")
}
