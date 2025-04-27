package models

import (
	"log"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/constants"
	"github.com/google/uuid"
)

func GenerateAndSaveMonthlyRequestCounts() error {
	var apps []App
	startTime := time.Now()
	if err := db.Find(&apps).Error; err != nil {
		return err
	}

	for _, app := range apps {
		log.Printf("=== App name %s ===", app.Name)

		var months []struct {
			YearMonth string
		}

		err := db.Model(&Request{}).
			Select("to_char(\"createdAt\", 'YYYY-MM') as year_month").
			Where("\"appId\" = ?", app.ID).
			Group("year_month").
			Scan(&months).Error
		if err != nil {
			log.Printf("Error selecting months: %+v", err)
			continue
		}

		log.Printf("months: %+v", months)

		for _, month := range months {
			// Skip the current month
			currentMonth := time.Now().Format("2006-01")
			if month.YearMonth == currentMonth {
				log.Printf("Skipping current month (%s) for AppID %s", month.YearMonth, app.ID)
				continue
			}
			var existingRC RequestCount

			// Check if RequestCount already exists for this App and Month
			if err := db.Where("\"appId\" = ? AND \"month\" = ?", app.ID, month.YearMonth).
				First(&existingRC).Error; err != nil && err.Error() != constants.RECORD_NOT_FOUND_ERROR {
				log.Printf("Error checking months: %+v", err)
				continue
			}
			if existingRC.ID != "" {
				log.Printf("RequestCount already exists for AppID %s and Month %s, skipping...", app.ID, month.YearMonth)
				continue
			}

			var firstRequest, lastRequest Request
			var requestCount int64

			startOfMonth, _ := time.Parse("2006-01", month.YearMonth)
			endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

			// First request in the month
			log.Printf("Getting first request: %+v", app.Name)
			if err := db.Where("\"appId\" = ? AND \"createdAt\" >= ? AND \"createdAt\" <= ?", app.ID, startOfMonth, endOfMonth).
				Order("\"createdAt\" ASC").
				First(&firstRequest).Error; err != nil {
				log.Printf("Error selecting first request: %+v", err)
				continue
			}

			// Last request in the month
			log.Printf("Getting last request: %+v", app.Name)
			if err := db.Where("\"appId\" = ? AND \"createdAt\" >= ? AND \"createdAt\" <= ?", app.ID, startOfMonth, endOfMonth).
				Order("\"createdAt\" DESC").
				First(&lastRequest).Error; err != nil {
				log.Printf("Error selecting last request: %+v", err)
				continue
			}

			// Total request count for the month
			log.Printf("Getting request count: %s  month: %s", app.Name, month)
			if err := db.Model(&Request{}).
				Where("\"appId\" = ? AND \"createdAt\" >= ? AND \"createdAt\" <= ?", app.ID, startOfMonth, endOfMonth).
				Count(&requestCount).Error; err != nil {
				log.Printf("Error generating app request count: %+v", err)
				continue
			}

			rc := RequestCount{
				ID:        uuid.NewString(),
				AppID:     app.ID,
				Month:     month.YearMonth,
				Count:     requestCount,
				FromDate:  firstRequest.CreatedAt,
				ToDate:    lastRequest.CreatedAt,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			log.Printf("Creating request count: %+v", app.Name)
			if err := db.Create(&rc).Error; err != nil {
				log.Printf("Error creating app request count: %+v", err)
				continue
			}
		}
	}

	log.Println("generateRequestCountDuration:", time.Since(startTime))
	return nil
}

func init() {
	log.Println("===Initiated GenerateAndSaveMonthlyRequestCounts===")
	// Run the first task after app start (after 2 minutes)
	go runInitialRequestCountGeneration()

	// Run the scheduled monthly task
	go runMonthlyRequestCountGeneration()
}

func runInitialRequestCountGeneration() {
	log.Println("Waiting 2 minutes to run initial request count generation...")
	time.Sleep(2 * time.Minute)

	log.Println("Running initial request count generation after 2 minutes")
	if err := GenerateAndSaveMonthlyRequestCounts(); err != nil {
		log.Printf("Error running GenerateAndSaveMonthlyRequestCounts after 2 minutes: %+v", err)
	}
}

func runMonthlyRequestCountGeneration() {
	for {
		now := time.Now()

		// Calculate next first day of the month at midnight
		nextMonth := now.AddDate(0, 1, 0)
		nextMidnight := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, now.Location())

		// Time until the next month's midnight
		durationUntilNextMidnight := time.Until(nextMidnight)
		log.Printf("Sleeping until next month's midnight: %s", nextMidnight.Format(time.RFC3339))

		time.Sleep(durationUntilNextMidnight)

		log.Println("Running monthly request count generation...")
		if err := GenerateAndSaveMonthlyRequestCounts(); err != nil {
			log.Printf("Error running monthly GenerateAndSaveMonthlyRequestCounts: %+v", err)
		}
	}
}
