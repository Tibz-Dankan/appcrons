package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var requestCache = RequestCache{}

func (r *Request) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (r *Request) Create(request Request) (Request, error) {
	result := db.Create(&request)

	if result.Error != nil {
		return request, result.Error
	}
	return request, nil
}

func (r *Request) FindOne(id string) (Request, error) {
	var request Request
	var err error

	if request, err = requestCache.Read(id); err != nil {
		return request, err
	}

	if request.ID != "" {
		return request, nil
	}
	db.First(&request, "id = ?", id)

	if err = requestCache.Write(request); err != nil {
		return request, err
	}

	return request, nil
}

func (r *Request) FindByApp(appId string, createdAtBeforeCursor time.Time) ([]Request, int64, error) {
	var requests []Request

	// 12-hour threshold
	threshold := time.Now().Add(-12 * time.Hour)

	query := db.Table("requests").Where("\"appId\" = ?", appId).Where("\"createdAt\" > ?", threshold)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if !createdAtBeforeCursor.IsZero() {
		query = query.Where("\"createdAt\" < ?", createdAtBeforeCursor)
	}

	if err := query.Order("\"createdAt\" desc").Limit(10).Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, count, nil
}

func (r *Request) FindAll() ([]Request, error) {
	var request []Request
	db.Find(&request)

	return request, nil
}

func (r *Request) Delete(id string) error {

	if err := db.Unscoped().Where("id = ?", id).Delete(&Request{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Request) DeleteByApp(appId string) error {

	if err := db.Unscoped().Where("\"appId\" = ?", appId).Delete(&Request{}).Error; err != nil {
		return err
	}

	return nil
}
