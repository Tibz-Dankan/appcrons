package routes

import (
	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/routes/app"
	"github.com/Tibz-Dankan/keep-active/internal/routes/auth"
	"github.com/Tibz-Dankan/keep-active/internal/routes/feedback"
	"github.com/Tibz-Dankan/keep-active/internal/routes/request"

	"github.com/gorilla/mux"
)

func AppRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.Logger)
	// App routes
	appRouter := router.PathPrefix("/api/v1/apps").Subrouter()
	appRouter.Use(middlewares.Auth)
	app.PostAppRoute(appRouter)
	app.UpdateAppRoute(appRouter)
	app.GetAppRoute(appRouter)
	app.GetAppByUserRoute(appRouter)
	app.GetAllAppsRoute(appRouter)
	app.DisableAppRoute(appRouter)
	app.EnableAppRoute(appRouter)
	app.DeleteAppRoute(appRouter)
	app.SearchAppsRoute(appRouter)

	// Request routes
	requestRouter := router.PathPrefix("/api/v1/requests").Subrouter()
	requestRouter.Use(middlewares.Auth)
	request.GetRequestByUserRoute(requestRouter)
	request.GetRequestRoute(requestRouter)
	request.GetLiveRequestsRoute(requestRouter)
	request.PostRequestTimeRoute(requestRouter)
	request.UpdateRequestTimeRoute(requestRouter)
	request.DeleteRequestTimeRoute(requestRouter)

	// Auth routes
	authRouter := router.PathPrefix("/api/v1/auth").Subrouter()
	auth.SignUpRoute(authRouter)
	auth.SignInRoute(authRouter)
	auth.ForgotPasswordRoute(authRouter)
	auth.ResetPasswordRoute(authRouter)
	// Authorized Auth routes
	authorizedAuthRouter := router.PathPrefix("/api/v1/auth").Subrouter()
	authorizedAuthRouter.Use(middlewares.Auth)
	auth.UpdateUserDetailsRoute(authorizedAuthRouter)
	auth.ChangePasswordRoute(authorizedAuthRouter)

	// Feedback Routes
	feedbackRouter := router.PathPrefix("/api/v1/feedback").Subrouter()
	feedbackRouter.Use(middlewares.Auth)
	feedback.PostFeedbackRoute(feedbackRouter)

	// Active route
	GetActiveRoute(router)

	return router
}
