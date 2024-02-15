package request

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getLiveRequests(w http.ResponseWriter, r *http.Request) {

	// TODO: To use server sent events, goroutines and channels to service request operations on the client side

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(response)
}

func GetLiveRequestsRoute(router *mux.Router) {
	router.HandleFunc("/get-live", getLiveRequests).Methods("GET")
}
