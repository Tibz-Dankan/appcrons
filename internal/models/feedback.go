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

func (f *Feedback) Update() error {
	db.Save(&f)

	return nil
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
func (f *Feedback) FindAll(createdAtBeforeCursor time.Time) ([]Feedback, int64, error) {
	var feedbacks, userFeedback []Feedback

	query := db.Table("feedbacks")

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if !createdAtBeforeCursor.IsZero() {
		query = query.Where("\"createdAt\" < ?", createdAtBeforeCursor)
	}

	if err := query.Order("\"createdAt\" desc").Limit(10).Find(&feedbacks).Error; err != nil {
		return nil, 0, err
	}

	for _, feedback := range feedbacks {
		var user User
		db.Select("id", "name", "email", "role", "\"createdAt\"", "\"updatedAt\"").First(&user, "id = ?", feedback.UserID)
		feedback.User = user

		userFeedback = append(userFeedback, feedback)
	}

	return userFeedback, count, nil
}
