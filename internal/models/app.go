package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (a *App) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	tx.Statement.SetColumn("IsDisabled", true)
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

// Fetches application details without
// including its requestTime and last request.
func (a *App) FindAppDetails(id string) (App, error) {
	var app App

	db.First(&app, "id = ?", id)

	return app, nil
}

// Fetches application details and includes its
// requestTime and last request.
func (a *App) FindOne(id string) (App, error) {
	var app App

	db.Preload("RequestTime").Preload("Request", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"createdAt\" desc").Limit(1)
	}).First(&app, "id = ?", id)

	return app, nil
}

func (a *App) FindByUser(userId string) ([]App, error) {
	var apps []App
	// TODO: To add cursor based pagination for fetching at maximum 10 apps for each user
	// TODO: To remove the last last request for each on this query

	startTime := time.Now()
	result := db.Preload("RequestTime").Order("\"updatedAt\" desc").Find(&apps, "\"userId\" = ?", userId)
	if result.Error != nil {
		return apps, nil
	}

	for i, app := range apps {
		var requests []Request
		db.Order("\"createdAt\" desc").Limit(1).Find(&requests, "\"appId\" = ?", app.ID)
		app.Request = requests
		apps[i] = app
	}

	log.Println("queryTimeMS:", int(time.Since(startTime).Milliseconds()))

	return apps, nil
}

// // Approach 0: Original naive implementation without optimization
// func (a *App) FindByUserPaginated(userId string, limit float64, cursor string) ([]App, error) {
// 	startTime := time.Now()
// 	var apps []App
	
// 	query := db.Model(&App{}).
// 	    Preload("RequestTime").
// 		Order("\"createdAt\" DESC").
// 		Limit(int(limit))

// 	if cursor != "" {
// 		var lastApp App
// 		if err := db.Select("\"createdAt\"").Where("id = ?",
// 			cursor).First(&lastApp).Error; err != nil {
// 			return apps, err
// 		}
// 		query = query.Where("\"createdAt\" < ?", lastApp.CreatedAt)
// 	}

// 	query.Find(&apps)

// 	// To be kept for comparison
// 	for i, app := range apps {
// 		var requests []Request
// 		db.Order("\"createdAt\" desc").Limit(1).Find(&requests, "\"appId\" = ?", app.ID)
// 		app.Request = requests
// 		apps[i] = app
// 	}

// 	log.Println("queryTimeMS:", int(time.Since(startTime).Milliseconds()))

// 	return apps, nil
// }

// Approach 1: Using Preload with conditions (GORM v2 feature)
func (a *App) FindByUserPaginated(userId string, limit float64, cursor string) ([]App, error) {
	startTime := time.Now()
	var apps []App
	var appCount int64

	countQuery := db.Model(&App{}).Order("\"createdAt\" DESC")
	if countQuery.Error != nil {
		return nil, countQuery.Error
	}

	if cursor != "" {
		var lastApp App
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).
		First(&lastApp).Error; err != nil {
			return apps, err
		}
		countQuery = countQuery.Where("\"createdAt\" < ?", lastApp.CreatedAt)
	}

	if err := countQuery.Count(&appCount).Error; err != nil {
		return apps, err
	}
	log.Println("Total apps for user in query:", appCount)
	
	query := db.Model(&App{}).
		Preload("RequestTime").
		Preload("Request", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"createdAt\" DESC").Limit(int(appCount))
		}).
		Order("\"createdAt\" DESC").
		Limit(int(limit))

	if cursor != "" {
		var lastApp App
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).First(&lastApp).Error; err != nil {
			return apps, err
		}
		query = query.Where("\"createdAt\" < ?", lastApp.CreatedAt)
	}

	if err := query.Find(&apps).Error; err != nil {
		return apps, err
	}

	log.Println("queryTimeMS:", int(time.Since(startTime).Milliseconds()))
	return apps, nil
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
	var apps []App

	log.Println("Fetching all apps")

	startTime := time.Now()
	result := db.Preload("RequestTime").Order("\"updatedAt\" desc").Find(&apps)
	if result.Error != nil {
		return apps, nil
	}

	// TODO: to find pagination solution for this part
	for i, app := range apps {
		var requests []Request
		db.Order("\"createdAt\" desc").Limit(1).Find(&requests, "\"appId\" = ?", app.ID)
		app.Request = requests
		apps[i] = app
	}

	log.Println("queryTimeMS:", int(time.Since(startTime).Milliseconds()))

	return apps, nil
}

// Fetches all user applications and includes
// each application's requestTimes.
func (a *App) FindByUserAndIncludeRequestTimes(userId string) ([]App, error) {
	var apps []App

	startTime := time.Now()
	result := db.Preload("RequestTime").Order("\"updatedAt\" desc").Find(&apps, "\"userId\" = ?", userId)
	if result.Error != nil {
		return apps, nil
	}

	log.Println("queryTimeMS:", int(time.Since(startTime).Milliseconds()))

	return apps, nil
}

// Search for apps that belong to the given userId
// and whose name contains the query string
func (a *App) Search(query, userId string) ([]App, error) {
	var apps []App

	result := db.Where("\"userId\" = ? AND name ILIKE ?", userId, "%"+query+"%").Order("\"updatedAt\" desc").Find(&apps)
	if result.Error != nil {
		return apps, result.Error
	}

	for i, app := range apps {
		var requests []Request
		db.Order("\"createdAt\" desc").Limit(1).Find(&requests, "\"appId\" = ?", app.ID)

		app.Request = requests
		apps[i] = app
	}

	return apps, nil
}

func (a *App) FindCount() (int64, error) {
	startTime := time.Now()
	var appCount int64

	if err := db.Model(&App{}).Count(&appCount).Error; err != nil {
		return appCount, err
	}
	log.Println("Total app count:", appCount)
	log.Println("queryTimeMS:", int(time.Since(startTime).Milliseconds()))

	return appCount, nil
}

// Update updates one app in the database, using the information
// stored in the receiver a
func (a *App) Update() (App, error) {
	db.Save(&a)

	app, err := a.FindOne(a.ID)
	if err != nil {
		return app, err
	}

	return app, err
}

func (a *App) Delete(id string) error {

	request := Request{AppID: id}
	requestTime := RequestTime{AppID: id}

	if err := requestTime.DeleteByApp(id); err != nil {
		return err
	}
	if err := request.DeleteByApp(id); err != nil {
		return err
	}
	if err := db.Unscoped().Where("id = ?", id).Delete(&App{}).Error; err != nil {
		return err
	}

	return nil
}
