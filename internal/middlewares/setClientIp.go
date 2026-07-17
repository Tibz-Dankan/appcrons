package middlewares

import (
	"context"
	"log"
	"net/http"
	"os"
)

const ClientIPKey ContextKey = "clientIP"

// SetClientIp resolves the real client IP and stores it in the request
// context under ClientIPKey, for use by session/site-visit tracking.
func SetClientIp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientIP string

		if os.Getenv("GO_ENV") == "production" {
			clientIP = r.Header.Get("CF-Connecting-IP")
			log.Println("IP Cloudflare header: ", clientIP)
		}
		if clientIP == "" {
			clientIP = r.Header.Get("X-Forwarded-For")
			log.Println("IP X-Forwarded-For: ", clientIP)
		}
		if clientIP == "" {
			clientIP = r.Header.Get("X-Real-IP")
			log.Println("IP X-Real-IP: ", clientIP)
		}
		if clientIP == "" {
			clientIP = r.RemoteAddr
			log.Println("IP RemoteAddr: ", clientIP)
		}

		log.Println("SetClientIp address: ", clientIP)
		log.Println("SetClientIp address: ", r.Header.Get("User-Agent"))

		ctx := context.WithValue(r.Context(), ClientIPKey, clientIP)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
