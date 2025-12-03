package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (rc *RequestCount) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (rc *RequestCount) Create(requestCount RequestCount) (RequestCount, error) {
	result := db.Create(&requestCount)

	if result.Error != nil {
		return requestCount, result.Error
	}
	return requestCount, nil
}

func (rc *RequestCount) Delete(id string) error {

	if err := db.Unscoped().Where("id = ?", id).Delete(&RequestCount{}).Error; err != nil {
		return err
	}

	return nil
}

func (rc *RequestCount) DeleteByApp(appId string) error {

	if err := db.Unscoped().Where("\"appId\" = ?", appId).Delete(&RequestCount{}).Error; err != nil {
		return err
	}

	return nil
}
