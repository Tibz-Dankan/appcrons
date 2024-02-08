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
	// App routes
	appRouter := router.PathPrefix("/api/v1/app").Subrouter()
	appRouter.Use(middlewares.Auth)
	app.PostAppRoute(appRouter)
	app.GetAppRoute(appRouter)
	app.GetAppByUserByUser(appRouter)

	// Auth routes
	authRouter := router.PathPrefix("/api/v1/auth").Subrouter()
	auth.SignUpRoute(authRouter)
	auth.SignInRoute(authRouter)
	auth.ForgotPasswordRoute(authRouter)
	auth.ResetPasswordRoute(authRouter)

	return router
}
