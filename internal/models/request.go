package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var requestCache = RequestCache{}

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

func (r *Request) FindByApp(appId string) ([]Request, error) {
	var requests []Request
	var err error

	if requests, err = requestCache.ReadByApp(appId); err != nil {
		return requests, err
	}

	if len(requests) != 0 {
		return requests, nil
	}

	db.Find(&requests, "\"appId\" = ?", appId)

	if err = requestCache.WriteByApp(appId, requests); err != nil {
		return requests, err
	}

	return requests, nil
}

func (r *Request) FindAll() ([]Request, error) {
	var request []Request
	db.Find(&request)

	return request, nil
}

func (r *Request) Delete(id string) error {
	db.Delete(&Request{}, id)
	if err := requestCache.Delete(id); err != nil {
		return err
	}

	return nil
}
