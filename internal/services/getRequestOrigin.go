package services

import (
	"net/http"
)

func GetRequestOrigin(r *http.Request) string {
	origin := r.Header.Get("Origin")
	if origin != "" {
		return origin
	}
	return r.Header.Get("Referer")
}
