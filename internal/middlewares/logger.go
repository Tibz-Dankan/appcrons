package middlewares

import (
	"log"
	"net/http"
	"time"
)

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

const (
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorReset  = "\033[0m"
)

func statusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return colorGreen
	case code >= 300 && code < 400:
		return colorYellow
	default:
		return colorRed
	}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		color := statusColor(rw.statusCode)
		log.Printf(
			"%s%d%s %s %s %s\n",
			color,
			rw.statusCode,
			colorReset,
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
