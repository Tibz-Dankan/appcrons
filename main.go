package main

import (
	"log"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/event"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/routes"
	"github.com/Tibz-Dankan/keep-active/internal/routes/request"
	"github.com/Tibz-Dankan/keep-active/internal/services"

	"github.com/rs/cors"
)

func main() {
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

	models.DBAutoMigrate()
	services.UpdateCacheOnBoot()

	log.Println("Starting http server up on 8000")
	go http.ListenAndServe(":8000", nil)

	services.StartClearUserAppMemoryScheduler()
	// Call StartRequestScheduler after server is started
	request.StartRequestScheduler()

	event.EventSubscribers()

	select {}
}
