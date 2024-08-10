package routes

import (
	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/routes/app"
	"github.com/Tibz-Dankan/keep-active/internal/routes/auth"
	"github.com/Tibz-Dankan/keep-active/internal/routes/feedback"
	"github.com/Tibz-Dankan/keep-active/internal/routes/monitor"
	"github.com/Tibz-Dankan/keep-active/internal/routes/request"

	"github.com/gorilla/mux"
)

func AppRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.RequestDuration)
	router.Use(middlewares.Logger)
	router.Use(middlewares.RateLimit)

	// App routes
	appRouter := router.PathPrefix("/api/v1/apps").Subrouter()
	appRouter.Use(middlewares.Auth)
	appRouter.Use(middlewares.HasPermissions)
	app.PostAppRoute(appRouter)
	app.UpdateAppRoute(appRouter)
	app.GetAppRoute(appRouter)
	app.GetAppByUserRoute(appRouter)
	app.GetAllAppsRoute(appRouter)
	app.DisableAppRoute(appRouter)
	app.EnableAppRoute(appRouter)
	app.DeleteAppRoute(appRouter)
	app.SearchAppsRoute(appRouter)
	app.GetAppsLastRequestByUserRoute(appRouter)

	// Request routes
	requestRouter := router.PathPrefix("/api/v1/requests").Subrouter()
	requestRouter.Use(middlewares.Auth)
	requestRouter.Use(middlewares.HasPermissions)
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
	authorizedAuthRouter.Use(middlewares.HasPermissions)
	auth.UpdateUserDetailsRoute(authorizedAuthRouter)
	auth.ChangePasswordRoute(authorizedAuthRouter)

	// Feedback Routes
	feedbackRouter := router.PathPrefix("/api/v1/feedback").Subrouter()
	feedbackRouter.Use(middlewares.Auth)
	feedbackRouter.Use(middlewares.HasPermissions)
	feedback.PostFeedbackRoute(feedbackRouter)
	feedback.GetFeedbackByUserRoute(feedbackRouter)
	feedback.GetAllFeedbackRoute(feedbackRouter)
	feedback.UpdateFeedbackRoute(feedbackRouter)

	// Active route
	GetActiveRoute(router)

	// Monitor route
	monitor.GetMetrics(router)

	return router
}
