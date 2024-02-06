package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                     string         `gorm:"column:id;type:uuid;primaryKey"`
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
	ID              string         `gorm:"column:id;type:uuid;primaryKey"`
	UserID          string         `gorm:"column:userId;not null;index"`
	Name            string         `gorm:"column:name;unique;not null;index"`
	URL             string         `gorm:"column:url;unique;not null"`
	RequestInterval string         `gorm:"column:requestInterval;not null"`
	Request         []Request      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RequestTime     []RequestTime  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt       time.Time      `gorm:"column:createdAt"`
	UpdatedAt       time.Time      `gorm:"column:updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deletedAt;index"`
}

type Request struct {
	ID        string         `gorm:"column:id;type:uuid;primaryKey"`
	AppID     string         `gorm:"column:appId;not null;index"`
	Status    int            `gorm:"column:status;not null"`
	Duration  string         `gorm:"column:duration;not null"`
	CreatedAt time.Time      `gorm:"column:createdAt"`
	UpdatedAt time.Time      `gorm:"column:updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deletedAt;index"`
}

type RequestTime struct {
	ID        string         `gorm:"column:id;type:uuid;primaryKey"`
	AppID     string         `gorm:"column:appId;not null;index"`
	Start     string         `gorm:"column:start;not null"`
	End       string         `gorm:"column:end;not null"`
	CreatedAt time.Time      `gorm:"column:createdAt"`
	UpdatedAt time.Time      `gorm:"column:updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deletedAt;index"`
}
