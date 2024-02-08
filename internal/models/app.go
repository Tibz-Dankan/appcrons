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

func (a *App) Create(app App) (string, error) {
	result := db.Create(&app)

	if result.Error != nil {
		return "", result.Error
	}
	return app.ID, nil
}

func (a *App) FindOne(id string) (App, error) {
	var app App
	db.First(&app, id)

	return app, nil
}

func (a *App) FindByUser(userId string) (App, error) {
	var app App
	db.Find(&app, "userId = ?", userId)

	return app, nil
}

func (a *App) FindByName(name string) (App, error) {
	var app App
	db.First(&app, "name = ?", name)

	return app, nil
}

func (a *App) FindAll() ([]App, error) {
	var apps []App
	db.Find(&apps)

	return apps, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (a *App) Update() error {
	db.Save(&a)

	return nil
}

func (a *App) Delete(id string) error {
	db.Delete(&App{}, id)

	return nil
}
