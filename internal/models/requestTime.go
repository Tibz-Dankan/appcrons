package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *RequestTime) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (r *RequestTime) Create(requestTime RequestTime) (RequestTime, error) {
	result := db.Create(&requestTime)

	if result.Error != nil {
		return requestTime, result.Error
	}
	return requestTime, nil
}

func (r *RequestTime) FindOne(id string) (RequestTime, error) {
	var request RequestTime
	db.First(&request, "id = ?", id)

	return request, nil
}

func (r *RequestTime) FindByApp(appId string) ([]RequestTime, error) {
	var requestTimes []RequestTime

	db.Find(&requestTimes, "\"appId\" = ?", appId)

	return requestTimes, nil
}

func (r *RequestTime) Delete(id string) error {
	db.Delete(&Request{}, id)

	return nil
}
