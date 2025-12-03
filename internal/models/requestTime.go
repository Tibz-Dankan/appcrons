package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (rc *RequestTime) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (rc *RequestTime) Create(requestTime RequestTime) (RequestTime, error) {
	result := db.Create(&requestTime)

	if result.Error != nil {
		return requestTime, result.Error
	}
	return requestTime, nil
}

func (rc *RequestTime) FindOne(id string) (RequestTime, error) {
	var requestTime RequestTime
	db.First(&requestTime, "id = ?", id)

	return requestTime, nil
}

func (rc *RequestTime) FindByApp(appId string) ([]RequestTime, error) {
	var requestTimes []RequestTime

	db.Find(&requestTimes, "\"appId\" = ?", appId)

	return requestTimes, nil
}

func (rc *RequestTime) Update() error {
	db.Save(&rc)

	return nil
}

func (rc *RequestTime) UpdateTimeZone(timeZone string) ([]RequestTime, error) {
	var requestTimes []RequestTime

	if err := db.Model(&rc).Where("\"appId\" = ?", rc.AppID).Update("timeZone", timeZone).Error; err != nil {
		return requestTimes, err
	}

	requestTimes, err := rc.FindByApp(rc.AppID)
	if err != nil {
		return requestTimes, err
	}

	return requestTimes, nil
}

func (r *RequestTime) Delete(id string) error {

	if err := db.Unscoped().Where("id = ?", id).Delete(&RequestTime{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *RequestTime) DeleteByApp(appId string) error {

	if err := db.Unscoped().Where("\"appId\" = ?", appId).Delete(&RequestTime{}).Error; err != nil {
		return err
	}

	return nil
}
