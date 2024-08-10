package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var requestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "go_requests_duration_seconds",
		Help:    "Duration of HTTP requests in seconds.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"path", "method", "status"},
)

// Initializes the registering of request duration
// by Prometheus client
func InitRequestDurationPromRegister() {
	prometheus.MustRegister(requestDuration)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Captures the request data such as path, method, statusCode, and
// duration and registers it with Prometheus client.
// This happens for every request.
func RequestDuration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(rw.statusCode)
		path := r.URL.Path
		method := r.Method

		requestDuration.WithLabelValues(path, method, statusCode).Observe(duration)
	})
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
