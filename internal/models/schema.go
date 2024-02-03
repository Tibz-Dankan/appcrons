package models

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type User struct {
	ID                     uuid.UUID      `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	Name                   string         `gorm:"column:name;not null"`
	Email                  string         `gorm:"column:email;unique;not null"`
	Password               string         `gorm:"column:password;not null"`
	PasswordResetToken     string         `gorm:"column:passwordResetToken;index"`
	PasswordResetExpiresAt time.Time      `gorm:"column:passwordResetExpiresAt"`
	Role                   string         `gorm:"column:role;default:'admin';not null"`
	App                    []App          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt              time.Time      `gorm:"column:createdAt"`
	UpdatedAt              time.Time      `gorm:"column:updatedAt"`
	DeletedAt              gorm.DeletedAt `gorm:"column:deletedAt;index"`
}

type App struct {
	ID              uuid.UUID      `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID          string         `gorm:"column:userId;not null"`
	Name            string         `gorm:"column:name;unique;not null"`
	URL             string         `gorm:"column:url;unique;not null"`
	RequestInterval string         `gorm:"column:requestInterval;not null"`
	Request         []Request      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt       time.Time      `gorm:"column:createdAt"`
	UpdatedAt       time.Time      `gorm:"column:updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deletedAt;index"`
}

type Request struct {
	ID        uuid.UUID      `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	AppID     string         `gorm:"column:appId;not null"`
	Status    int            `gorm:"column:status;not null"`
	CreatedAt time.Time      `gorm:"column:createdAt"`
	UpdatedAt time.Time      `gorm:"column:updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deletedAt;index"`
}
