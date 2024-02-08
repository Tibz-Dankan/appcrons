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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "production_url"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
		Debug:            true,
		AllowedHeaders:   []string{"*"},
	})

	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file=> ", err.Error())
	// }

	handler := c.Handler(router)

	http.Handle("/", handler)

	models.DBAutoMigrate()

	fmt.Println("Starting http server up on 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
