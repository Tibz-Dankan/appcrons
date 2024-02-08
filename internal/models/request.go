package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Request) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (r *Request) Create(request Request) (string, error) {
	result := db.Create(&request)

	if result.Error != nil {
		return "", result.Error
	}
	return request.ID, nil
}

func (r *Request) FindOne(id string) (Request, error) {
	var request Request
	db.First(&request, id)

	return request, nil
}

func (r *Request) FindByApp(appId string) (Request, error) {
	var request Request
	db.Find(&request, "appId = ?", appId)

	return request, nil
}

func (r *Request) FindAll() ([]Request, error) {
	var request []Request
	db.Find(&request)

	return request, nil
}

func (r *Request) Delete(id string) error {
	db.Delete(&Request{}, id)

	return nil
}
