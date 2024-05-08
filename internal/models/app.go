package models

import (
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

	if err := db.First(&app, "id = ?", app.ID).Error; err != nil {
		return app, err
	}

	return app, nil
}

func (a *App) FindOne(id string) (App, error) {
	var app App
	var err error

	if app, err = appCache.Read(id); err != nil {
		return app, err
	}

	if app.ID != "" {
		return app, nil
	}
	db.First(&app, "id = ?", id)

	if err = appCache.Write(app); err != nil {
		return app, err
	}

	return app, nil
}

func (a *App) FindByUser(userId string) ([]App, error) {
	var apps []App
	var err error

	if apps, err = appCache.ReadByUser(userId); err != nil {
		return apps, err
	}

	if len(apps) != 0 {
		return apps, nil
	}

	db.Find(&apps, "\"userId\" = ?", userId)

	if err = appCache.WriteByUser(userId, apps); err != nil {
		return apps, err
	}

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
	db.Find(&apps)
	// TODO: To add redis read and write

	return apps, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (a *App) Update() error {
	db.Save(&a)

	if err := appCache.Write(*a); err != nil {
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
