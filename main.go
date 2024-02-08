package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/routes"

	"github.com/rs/cors"
)

func main() {
	router := routes.AppRouter()
	// TODO: to add rate limiting

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

	fmt.Println("Starting http server up on 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
