package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (l *Location) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (l *Location) Create(location Location) (Location, error) {
	result := db.Create(&location)

	if result.Error != nil {
		return location, result.Error
	}

	return location, nil
}

func (l *Location) FindOne(id string) (Location, error) {
	var location Location

	db.First(&location, "id = ?", id)

	return location, nil
}

// FindByIP finds the most recently created Location for a given IP,
// used to avoid re-resolving geo-IP info for an IP already seen before.
func (l *Location) FindByIP(ip string) (Location, error) {
	var location Location

	db.Where("ip = ?", ip).Order("\"createdAt\" desc").First(&location)

	return location, nil
}
