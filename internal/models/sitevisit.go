package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *SiteVisit) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (s *SiteVisit) Create(siteVisit SiteVisit) (SiteVisit, error) {
	result := db.Create(&siteVisit)

	if result.Error != nil {
		return siteVisit, result.Error
	}

	return siteVisit, nil
}

func (s *SiteVisit) FindOne(id string) (SiteVisit, error) {
	var siteVisit SiteVisit

	db.First(&siteVisit, "id = ?", id)

	return siteVisit, nil
}
