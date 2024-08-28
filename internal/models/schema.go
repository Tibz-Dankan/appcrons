package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var db = Db()
var DB = db
var redisClient = RedisClient()
var ctx = context.Background()

type User struct {
	ID                     string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name                   string         `gorm:"column:name;not null;index" json:"name"`
	Email                  string         `gorm:"column:email;unique;not null;index" json:"email"`
	Password               string         `gorm:"column:password;not null" json:"password"`
	PasswordResetToken     string         `gorm:"column:passwordResetToken;index" json:"passwordResetToken"`
	PasswordResetExpiresAt time.Time      `gorm:"column:passwordResetExpiresAt;index" json:"passwordResetExpiresAt"`
	Role                   string         `gorm:"column:role;default:'user';not null" json:"role"`
	App                    []App          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"apps"`
	Feedback               []Feedback     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"feedbacks"`
	CreatedAt              time.Time      `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt              time.Time      `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt              gorm.DeletedAt `gorm:"column:deletedAt;index" json:"deletedAt"`
}

type App struct {
	ID              string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID          string         `gorm:"column:userId;not null;index" json:"userId"`
	Name            string         `gorm:"column:name;unique;not null;index" json:"name"`
	URL             string         `gorm:"column:url;unique;not null;index" json:"url"`
	RequestInterval string         `gorm:"column:requestInterval;not null" json:"requestInterval"`
	IsDisabled      bool           `gorm:"column:isDisabled" json:"isDisabled"`
	Request         []Request      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"requests"`
	RequestTime     []RequestTime  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"requestTimes"`
	CreatedAt       time.Time      `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"column:updatedAt;index" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deletedAt;index" json:"deletedAt"`
}

type Request struct {
	ID         string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AppID      string         `gorm:"column:appId;not null;index" json:"appId"`
	StatusCode int            `gorm:"column:statusCode;not null" json:"statusCode"`
	Duration   int            `gorm:"column:duration;not null" json:"duration"`
	StartedAt  time.Time      `gorm:"column:startedAt;default:CURRENT_TIMESTAMP;index" json:"startedAt"`
	CreatedAt  time.Time      `gorm:"column:createdAt;index" json:"createdAt"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deletedAt;index" json:"deletedAt"`
}

type RequestTime struct {
	ID        string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AppID     string         `gorm:"column:appId;not null;index" json:"appId"`
	Start     string         `gorm:"column:start;not null" json:"start"`
	End       string         `gorm:"column:end;not null" json:"end"`
	TimeZone  string         `gorm:"column:timeZone;not null" json:"timeZone"`
	CreatedAt time.Time      `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updatedAt;index" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deletedAt;index" json:"deletedAt"`
}

type Feedback struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"column:userId;not null;index" json:"userId"`
	User      User      `gorm:"foreignKey:UserID;references:ID;default:null" json:"user"`
	Rating    float32   `gorm:"column:rating;not null" json:"rating"`
	Message   string    `gorm:"column:message;not null" json:"message"`
	CreatedAt time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}
