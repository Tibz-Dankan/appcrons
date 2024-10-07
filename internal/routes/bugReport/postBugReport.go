package bugreport

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func postBugReport(w http.ResponseWriter, r *http.Request) {
	bugReport := models.BugReport{}

	// Parse the multipart form with a 10 MB limit
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		services.AppError("Unable to parse form", 400, w)
		return
	}

	// Extract 'title' and 'description' values from the form
	bugReport.Title = r.FormValue("title")
	bugReport.Description = r.FormValue("description")

	log.Printf("bugReport: +%v ", bugReport)

	if bugReport.Title == "" || bugReport.Description == "" {
		services.AppError("Missing title or description!", 400, w)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}
	defer file.Close()

	log.Println("About to generate random number")

	randNumStr := strconv.Itoa(rand.Intn(9000) + 1000)
	filePath := "bugreports/" + randNumStr + "_" + fileHeader.Filename

	upload := services.Upload{FilePath: filePath}

	imageUrl, err := upload.Add(file, fileHeader)
	if err != nil {
		log.Println("err :", err)
		services.AppError(err.Error(), 500, w)
		return
	}

	bugReport.Image = imageUrl

	fmt.Println("imageUrl :", imageUrl)

	// TODO: upload image to firebase storage or aws s3 here

	newBugReport, err := bugReport.Create(bugReport)
	if err != nil {
		services.AppError(err.Error(), 400, w)
		return
	}

	response := map[string]interface{}{
		"status":    "success",
		"message":   "Thank you very much for reporting this bug",
		"bugReport": newBugReport,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func PostReportBugRoute(router *mux.Router) {
	router.HandleFunc("/post", postBugReport).Methods("POST")
}
