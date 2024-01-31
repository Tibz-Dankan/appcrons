package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                     int            `gorm:"column:id;primaryKey;autoIncrement"`
	Name                   string         `gorm:"column:name;not null"`
	Email                  string         `gorm:"column:email;unique;not null"`
	Password               string         `gorm:"column:password;not null"`
	PasswordResetToken     string         `gorm:"column:passwordResetToken;index"`
	PasswordResetExpiresAt time.Time      `gorm:"column:passwordResetExpiresAt"`
	Role                   string         `gorm:"column:role;default:'admin';not null"`
	CreatedAt              time.Time      `gorm:"column:createdAt"`
	UpdatedAt              time.Time      `gorm:"column:updatedAt"`
	DeletedAt              gorm.DeletedAt `gorm:"column:deletedAt;index"`
}

type App struct {
	ID              int            `gorm:"column:id;primaryKey;autoIncrement"`
	Name            string         `gorm:"column:name;not null"`
	URL             string         `gorm:"column:url;not null"`
	RequestInterval string         `gorm:"column:requestInterval;not null"`
	CreatedAt       time.Time      `gorm:"column:createdAt"`
	UpdatedAt       time.Time      `gorm:"column:updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deletedAt;index"`
}
