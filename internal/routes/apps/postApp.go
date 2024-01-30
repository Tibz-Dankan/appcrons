package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func PostAdd(w http.ResponseWriter, r *http.Request) {

	response := map[string]interface{}{
		"status":  "success",
		"message": "Created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func PostAppRoute(router *mux.Router) {
	router.HandleFunc("/api/v1/app/post", PostAdd).Methods("POST")
}
