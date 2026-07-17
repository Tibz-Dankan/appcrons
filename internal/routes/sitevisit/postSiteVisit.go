package sitevisit

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

type postSiteVisitInput struct {
	Page       string `json:"page"`
	Path       string `json:"path"`
	CapturedAt string `json:"capturedAt"`
}

func postSiteVisit(w http.ResponseWriter, r *http.Request) {
	input := postSiteVisitInput{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	if input.Page == "" || input.Path == "" || input.CapturedAt == "" {
		services.AppError("page, path and capturedAt are required!", 400, w)
		return
	}

	capturedAt, err := time.Parse(time.RFC3339, input.CapturedAt)
	if err != nil {
		services.AppError("Invalid capturedAt format! Must be an ISO 8601 string.", 400, w)
		return
	}

	clientIP, _ := r.Context().Value(middlewares.ClientIPKey).(string)
	userId, _ := r.Context().Value(middlewares.UserIDKey).(string)
	device := r.Header.Get("x-device")

	log.Println("postSiteVisit clientIP address: ", clientIP)
	log.Println("postSiteVisit userId address: ", userId)
	log.Println("postSiteVisit device address: ", device)

	location, err := services.GetUserLocationByIP(userId, clientIP)
	if err != nil {
		log.Println("Error resolving location for site visit:", err)
	}

	log.Printf("postSiteVisit location: %+v", location)

	siteVisit := models.SiteVisit{
		UserID:     userId,
		Device:     device,
		Page:       input.Page,
		Path:       input.Path,
		LocationID: location.ID,
		CapturedAt: capturedAt,
	}

	newSiteVisit, err := siteVisit.Create(siteVisit)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	response := map[string]interface{}{
		"status":    "success",
		"message":   "Visit captured successfully",
		"siteVisit": newSiteVisit,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func PostSiteVisitRoute(router *mux.Router) {
	router.HandleFunc("/post", postSiteVisit).Methods("POST")
}
