package services

import (
	"log"
	"sync"
	"time"
)

const MaintenanceMessage = "Appcrons is undergoing a maintenance, try again later"

// Hardcoded window for the Render -> Neon Postgres migration.
const (
	maintenanceWindowStartStr = "2026-08-26T22:00:00+03:00"
	maintenanceWindowEndStr   = "2026-08-27T02:00:00+03:00"
)

var (
	maintenanceWindowStart time.Time
	maintenanceWindowEnd   time.Time
	maintenanceWindowSet   bool
	maintenanceOnce        sync.Once
)

func loadMaintenanceWindow() {
	start, err := time.Parse(time.RFC3339, maintenanceWindowStartStr)
	if err != nil {
		log.Println("Invalid MAINTENANCE_WINDOW_START:", err)
		return
	}

	end, err := time.Parse(time.RFC3339, maintenanceWindowEndStr)
	if err != nil {
		log.Println("Invalid MAINTENANCE_WINDOW_END:", err)
		return
	}

	maintenanceWindowStart = start
	maintenanceWindowEnd = end
	maintenanceWindowSet = true
}

func IsMaintenanceActive() bool {
	maintenanceOnce.Do(loadMaintenanceWindow)

	if !maintenanceWindowSet {
		return false
	}

	now := time.Now()
	return now.Equal(maintenanceWindowStart) ||
		now.Equal(maintenanceWindowEnd) ||
		(now.After(maintenanceWindowStart) && now.Before(maintenanceWindowEnd))
}
