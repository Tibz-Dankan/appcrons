package middlewares

import (
	"log"
	"net/http"
	"sync"
	"time"

	services "github.com/Tibz-Dankan/keep-active/internal/services"
)

type RateLimiter struct {
	requests map[string]int
	limit    int
	window   time.Duration
	mutex    sync.Mutex
}

var rateLimiter = &RateLimiter{
	requests: make(map[string]int),
	window:   60 * time.Second, // 1 min
	limit:    20,
	mutex:    sync.Mutex{},
}

func (rl *RateLimiter) resetCount(clientIp string) {
	time.Sleep(rl.window)
	rl.mutex.Lock()
	delete(rl.requests, clientIp)
	rl.mutex.Unlock()
}

func (rl *RateLimiter) AllowRequest(clientIp string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Check the current count for the client
	count, found := rl.requests[clientIp]
	if !found || count < rl.limit {
		// Reset the count for a new window or increment for the current window
		if !found {
			go rl.resetCount(clientIp)
		}
		rl.requests[clientIp]++
		return true
	}

	return false
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.Header.Get("X-Real-IP")
		}
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}

		log.Println("clientIP address: ", clientIP)

		if !rateLimiter.AllowRequest(clientIP) {
			services.AppError("You have made too many requests!, Please try again later.", 429, w)
			log.Println("Request denied for:", clientIP)
			return
		}

		next.ServeHTTP(w, r)
	})
}
