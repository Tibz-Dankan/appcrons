package routes

import (
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	message := "Endpoint " + "'" + r.URL.Path + "'" + " does not exist!"
	services.AppError(message, 404, w)
}

func NotFoundRoute(router *mux.Router) {
	router.NotFoundHandler = http.HandlerFunc(notFound)
}
