package main

import (
	"fmt"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/routes"

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

	fmt.Println("Starting http server up on 3000")
	go http.ListenAndServe(":3000", nil)

	// Call StartRequestScheduler after server is started
	// request.StartRequestScheduler()

	select {}
}
