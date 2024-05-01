package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                     string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name                   string         `gorm:"column:name;not null" json:"name"`
	Email                  string         `gorm:"column:email;unique;not null" json:"email"`
	Password               string         `gorm:"column:password;not null" json:"password"`
	PasswordResetToken     string         `gorm:"column:passwordResetToken;index" json:"passwordResetToken"`
	PasswordResetExpiresAt time.Time      `gorm:"column:passwordResetExpiresAt" json:"passwordResetExpiresAt"`
	Role                   string         `gorm:"column:role;default:'admin';not null" json:"role"`
	App                    []App          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt              time.Time      `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt              time.Time      `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt              gorm.DeletedAt `gorm:"column:deletedAt;index"`
}

type App struct {
	ID              string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID          string         `gorm:"column:userId;not null;index" json:"userId"`
	Name            string         `gorm:"column:name;unique;not null;index" json:"name"`
	URL             string         `gorm:"column:url;unique;not null;index" json:"url"`
	RequestInterval string         `gorm:"column:requestInterval;not null" json:"requestInterval"`
	Request         []Request      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"request"`
	RequestTime     []RequestTime  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"requestTime"`
	CreatedAt       time.Time      `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deletedAt;index" json:"deletedAt"`
}

type Request struct {
	ID         string         `gorm:"column:id;type:uuid;primaryKey"`
	AppID      string         `gorm:"column:appId;not null;index"`
	StatusCode int            `gorm:"column:statusCode;not null"`
	Duration   int            `gorm:"column:duration;not null"`
	CreatedAt  time.Time      `gorm:"column:createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deletedAt;index"`
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
