package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (a *App) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
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
	db.First(&app, "id = ?", id)
	// TODO: To add redis read and write

	return app, nil
}

func (a *App) FindByUser(userId string) ([]App, error) {
	var apps []App
	db.Find(&apps, "\"userId\" = ?", userId)
	// TODO: To add redis read and write

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
	// TODO: To add redis write through update

	return nil
}

func (a *App) Delete(id string) error {
	db.Delete(&App{}, id)
	// TODO: To add redis delete through delete

	return nil
}
