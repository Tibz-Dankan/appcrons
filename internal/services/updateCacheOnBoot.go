package services

import (
	"log"

	"github.com/Tibz-Dankan/keep-active/internal/config"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"gorm.io/gorm"
)

// This fn must be called at the start of
// the server to update the cache with the
// latest data from primary db
func UpdateCacheOnBoot() {
	log.Println("Updating cache onBootUp")

	appCache := models.AppCache{}
	apps := []models.App{}
	var err error
	var db = config.Db()

	result := db.Preload("RequestTime").Preload("Request", func(db *gorm.DB) *gorm.DB {
		subQuery := db.Table("requests").
			Select("MAX(\"requests\".\"createdAt\")").
			Where("\"requests\".\"appId\" = apps.id").
			Group("\"requests\".\"appId\"")
		return db.Where("\"requests\".\"createdAt\" IN (?)", subQuery).Joins("JOIN apps ON apps.id = \"requests\".\"appId\"")
	}).Find(&apps)
	if result.Error != nil {
		log.Println("Error writing to cache onBootUp: ", err)
		return
	}

	if err := appCache.WriteAll(apps); err != nil {
		log.Println("Error writing to cache onBootUp: ", err)
		return
	}
}
