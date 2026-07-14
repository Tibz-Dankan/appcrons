package middlewares

import (
	"net/http"

	services "github.com/Tibz-Dankan/keep-active/internal/services"
)

func Maintenance(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/get/active" || r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		if services.IsMaintenanceActive() {
			w.Header().Set("X-Maintenance-Mode", "true")
			services.AppError(services.MaintenanceMessage, http.StatusServiceUnavailable, w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
