package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var appCache = AppCache{}

func (a *App) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	tx.Statement.SetColumn("IsDisabled", false)
	return nil
}

func (a *App) Create(app App) (App, error) {
	result := db.Create(&app)

	if result.Error != nil {
		return app, result.Error
	}

	app, err := a.FindOne(app.ID)
	if err != nil {
		return app, err
	}

	return app, nil
}

func (a *App) FindOne(id string) (App, error) {
	var app App

	db.Preload("RequestTime").Preload("Request", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"createdAt\" DESC").Limit(1)
	}).First(&app, "id = ?", id)

	return app, nil
}

func (a *App) FindByUser(userId string) ([]App, error) {
	var apps, userApps []App
	var err error

	if apps, err = appCache.ReadByUser(userId); err != nil {
		return apps, err
	}

	if len(apps) != 0 {
		return apps, nil
	}

	startTime := time.Now()
	result := db.Preload("RequestTime").Order("\"updatedAt\" desc").Find(&apps, "\"userId\" = ?", userId)
	if result.Error != nil {
		return apps, nil
	}

	for _, app := range apps {
		var requests []Request
		db.Order("\"createdAt\" desc").Limit(1).Find(&requests, "\"appId\" = ?", app.ID)
		app.Request = requests

		userApps = append(userApps, app)
	}

	duration := time.Since(startTime)
	queryTimeMS := int(duration.Milliseconds())
	log.Println("queryTimeMS:", queryTimeMS)

	return userApps, nil
}

func (a *App) FindByName(name string) (App, error) {
	var app App
	db.First(&app, "name = ?", name)
	// TODO: To add redis read and write

	return app, nil
}

func (a *App) FindByURL(url string) (App, error) {
	var app App
	db.First(&app, "url = ?", url)
	// TODO: To add redis read and write

	return app, nil
}

func (a *App) FindAll() ([]App, error) {
	var apps, savedApps []App
	var err error

	if apps, err = appCache.ReadAll(); err != nil {
		return apps, err
	}

	if len(apps) != 0 {
		return apps, nil
	}

	log.Println("Fetching all apps")

	startTime := time.Now()
	result := db.Preload("RequestTime").Order("\"updatedAt\" desc").Find(&apps)
	if result.Error != nil {
		return apps, nil
	}

	// TODO: to find pagination solution for this part
	for _, app := range apps {
		var requests []Request
		db.Order("\"createdAt\" DESC").Limit(1).Find(&requests, "\"appId\" = ?", app.ID)
		app.Request = requests

		savedApps = append(savedApps, app)
	}

	duration := time.Since(startTime)
	queryTimeMS := int(duration.Milliseconds())
	log.Println("queryTimeMS:", queryTimeMS)

	if err = appCache.WriteAll(savedApps); err != nil {
		return apps, err
	}

	return savedApps, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (a *App) Update() error {
	db.Save(&a)

	if _, err := a.FindOne(a.ID); err != nil {
		return err
	}

	return nil
}

func (a *App) Delete(id string) error {
	db.Delete(&App{}, id)

	if err := appCache.Delete(id); err != nil {
		return err
	}

	return nil
}
