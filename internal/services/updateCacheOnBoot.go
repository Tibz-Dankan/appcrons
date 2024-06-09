package services

import (
	"log"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/config"
	"github.com/Tibz-Dankan/keep-active/internal/models"
)

// This fn must be called at the start of
// the server to update the cache with the
// latest data from primary db
func UpdateCacheOnBoot() {
	log.Println("Updating cache onBootUp")

	appCache := models.AppCache{}
	apps := []models.App{}
	savedApps := []models.App{}

	var err error
	var db = config.Db()

	startTime := time.Now()
	result := db.Preload("RequestTime").Order("\"updatedAt\" desc").Find(&apps)
	if result.Error != nil {
		log.Println("Error writing to cache onBootUp: ", err)
		return
	}

	for _, app := range apps {
		var requests []models.Request
		db.Order("\"createdAt\" desc").Limit(1).Find(&requests, "\"appId\" = ?", app.ID)
		app.Request = requests

		savedApps = append(savedApps, app)
	}

	duration := time.Since(startTime)
	queryTimeMS := int(duration.Milliseconds())

	log.Println("queryTimeMS:", queryTimeMS)

	if err := appCache.WriteAll(savedApps); err != nil {
		log.Println("Error writing to cache onBootUp: ", err)
		return
	}
	log.Println("Cache updated successfully")
}
