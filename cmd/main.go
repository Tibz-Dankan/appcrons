package main

import (
	"log"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/events/publishers"
	"github.com/Tibz-Dankan/keep-active/internal/events/subscribers"
	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/routes"
	"github.com/Tibz-Dankan/keep-active/internal/routes/auth"
	"github.com/Tibz-Dankan/keep-active/internal/schedulers"

	"github.com/rs/cors"
)

func main() {
	middlewares.InitRequestDurationPromRegister()
	router := routes.AppRouter()

	// Completes the migration models.Db() intentionally left unfinished:
	// backfills legacy Location/SiteVisit rows and creates their userId/
	// locationId foreign key constraints (see internal/models/db.go for
	// why this can't run during package initialization). Must run before
	// the server starts accepting traffic.
	models.FinishMigration()
	auth.CreateUnknownUser()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	handler := c.Handler(router)

	http.Handle("/", handler)

	log.Println("Starting http server up on 8080")
	go http.ListenAndServe(":8080", nil)

	go schedulers.InitSchedulers()
	go subscribers.InitEventSubscribers()
	publishers.InitEventPublishers()

	select {}
}
