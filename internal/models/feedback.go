package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (f *Feedback) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (f *Feedback) Create(feedback Feedback) (Feedback, error) {
	result := db.Create(&feedback)

	if result.Error != nil {
		return feedback, result.Error
	}

	feedback, err := f.FindOne(feedback.ID)
	if err != nil {
		return feedback, err
	}

	return feedback, nil
}

// Gets feedback by id
func (f *Feedback) FindOne(id string) (Feedback, error) {
	var feedback Feedback

	db.First(&feedback, "id = ?", id)

	return feedback, nil
}

// Gets paginated feedback for a specific user
func (f *Feedback) FindByUser(userId string, createdAtBeforeCursor time.Time) ([]Feedback, int64, error) {
	var feedback []Feedback

	query := db.Table("feedbacks").Where("\"userId\" = ?", userId)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return feedback, 0, err
	}

	if !createdAtBeforeCursor.IsZero() {
		query = query.Where("\"createdAt\" < ?", createdAtBeforeCursor)
	}

	if err := query.Order("\"createdAt\" desc").Limit(10).Find(&feedback).Error; err != nil {
		return feedback, 0, err
	}

	return feedback, count, nil
}

// Gets paginated feedback for all users
func (r *Feedback) FindAll(createdAtBeforeCursor time.Time) ([]Feedback, int64, error) {
	var feedback []Feedback

	query := db.Table("feedbacks").Preload("users", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email", "role", "\"createdAt\"", "\"createdAt\"")
	})

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if !createdAtBeforeCursor.IsZero() {
		query = query.Where("\"createdAt\" < ?", createdAtBeforeCursor)
	}

	if err := query.Order("\"createdAt\" desc").Limit(10).Find(&feedback).Error; err != nil {
		return nil, 0, err
	}

	return feedback, count, nil
}
