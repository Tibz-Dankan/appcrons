package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (s *Session) Create(session Session) (Session, error) {
	result := db.Create(&session)

	if result.Error != nil {
		return session, result.Error
	}

	session, err := s.FindOne(session.ID)
	if err != nil {
		return session, err
	}

	return session, nil
}

func (s *Session) FindOne(id string) (Session, error) {
	var session Session

	db.First(&session, "id = ?", id)

	return session, nil
}

// FindByUser gets paginated sessions for a specific user, most recent first.
func (s *Session) FindByUser(userId string, limit int) ([]Session, error) {
	var sessions []Session

	if err := db.Where("\"userId\" = ?", userId).
		Order("\"createdAt\" desc").Limit(limit).Find(&sessions).Error; err != nil {
		return sessions, err
	}

	return sessions, nil
}
