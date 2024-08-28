package main

import (
	"log"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/events/publishers"
	"github.com/Tibz-Dankan/keep-active/internal/events/subscribers"
	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/routes"
	"github.com/Tibz-Dankan/keep-active/internal/schedulers"

	"github.com/rs/cors"
)

func main() {
	middlewares.InitRequestDurationPromRegister()
	router := routes.AppRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
		Debug:            true,
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
