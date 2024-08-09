package monitor

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GetMetrics(router *mux.Router) {
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
}
