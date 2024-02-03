package routes

import (
	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/routes/app"

	"github.com/Tibz-Dankan/keep-active/internal/routes/auth"

	"github.com/gorilla/mux"
)

func AppRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.Logger)

	// Auth routes
	auth.SignUpRoute(router)
	auth.SignInRoute(router)
	auth.ForgotPasswordRoute(router)
	auth.ResetPasswordRoute(router)

	router.Use(middlewares.Auth)
	app.PostAppRoute(router)
	// App routes

	return router
}
