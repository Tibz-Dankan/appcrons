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
	Password               string         `gorm:"column:password;not null" json:"password,omitempty"`
	PasswordResetToken     string         `gorm:"column:passwordResetToken;index" json:"passwordResetToken,omitempty"`
	PasswordResetExpiresAt time.Time      `gorm:"column:passwordResetExpiresAt;index" json:"passwordResetExpiresAt,omitempty"`
	Role                   string         `gorm:"column:role;default:'user';not null" json:"role"`
	App                    []App          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"apps,omitempty"`
	Feedback               []Feedback     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"feedbacks,omitempty"`
	OPT                    []OTP          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"OPT,omitempty"`
	CreatedAt              time.Time      `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt              time.Time      `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt              gorm.DeletedAt `gorm:"column:deletedAt;index" json:"deletedAt,omitempty"`
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
	RequestCount    []RequestCount `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"requestCount"`
	CreatedAt       time.Time      `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"column:updatedAt;index" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deletedAt;index" json:"deletedAt,omitempty"`
}

type Request struct {
	ID         string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AppID      string         `gorm:"column:appId;not null;index" json:"appId"`
	StatusCode int            `gorm:"column:statusCode;not null" json:"statusCode"`
	Duration   int            `gorm:"column:duration;not null" json:"duration"`
	StartedAt  time.Time      `gorm:"column:startedAt;default:CURRENT_TIMESTAMP;index" json:"startedAt"`
	CreatedAt  time.Time      `gorm:"column:createdAt;index" json:"createdAt"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deletedAt;index" json:"deletedAt,omitempty"`
}

type RequestTime struct {
	ID        string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AppID     string         `gorm:"column:appId;not null;index" json:"appId"`
	Start     string         `gorm:"column:start;not null" json:"start"`
	End       string         `gorm:"column:end;not null" json:"end"`
	TimeZone  string         `gorm:"column:timeZone;not null" json:"timeZone"`
	CreatedAt time.Time      `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updatedAt;index" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deletedAt;index" json:"deletedAt,omitempty"`
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

type OTP struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"column:userId;not null;index" json:"userId"`
	OTP       int       `gorm:"column:OTP;index" json:"OTP"`
	IsUsed    bool      `gorm:"column:isUsed;default:false" json:"isUsed"`
	ExpiresAt time.Time `gorm:"column:expiresAt;not null;index" json:"expiresAt"`
	CreatedAt time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type BugReport struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      string    `gorm:"column:userId;default:null;index" json:"userId"`
	User        User      `gorm:"foreignKey:UserID;references:ID;default:null" json:"user"`
	Title       string    `gorm:"column:title;not null" json:"title"`
	Description string    `gorm:"column:description;not null" json:"description"`
	Image       string    `gorm:"column:image;default:null" json:"image"`
	CreatedAt   time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type RequestCount struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AppID     string    `gorm:"column:appId;not null;index" json:"appId"`
	Month     string    `gorm:"column:month;not null;index" json:"month"` // Format: "YYYY-MM"
	Count     int64     `gorm:"column:count;not null" json:"count"`
	FromDate  time.Time `gorm:"column:fromDate;not null" json:"fromDate"`
	ToDate    time.Time `gorm:"column:toDate;not null" json:"toDate"`
	CreatedAt time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type Location struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"column:userId;default:null;index" json:"userId"`
	IP        string    `gorm:"column:ip;not null;index" json:"ip"`
	Info      string    `gorm:"column:info;type:text;not null" json:"info"`
	CreatedAt time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	User      User      `gorm:"foreignKey:UserID;references:ID;default:null" json:"user,omitempty"`
}

type Session struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID       string    `gorm:"column:userId;not null;index" json:"userId"`
	AccessToken  string    `gorm:"column:accessToken;not null;index" json:"accessToken"`
	GeneratedVia string    `gorm:"column:generatedVia;not null;index" json:"generatedVia"`
	Device       string    `gorm:"column:device;type:text;default:'Unknown Device'" json:"device"`
	LocationID   string    `gorm:"column:locationId;default:null;index" json:"locationId"`
	IsRevoked    bool      `gorm:"column:isRevoked;default:false" json:"isRevoked"`
	CreatedAt    time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	User         User      `gorm:"foreignKey:UserID;references:ID;default:null" json:"user,omitempty"`
	Location     Location  `gorm:"foreignKey:LocationID;references:ID;default:null" json:"location,omitempty"`
}

type SiteVisit struct {
	ID         string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID     string    `gorm:"column:userId;default:null;index" json:"userId"`
	Device     string    `gorm:"column:device;type:text;default:'Unknown Device'" json:"device"`
	Page       string    `gorm:"column:page;not null;index" json:"page"`
	Path       string    `gorm:"column:path;not null;index" json:"path"`
	LocationID string    `gorm:"column:locationId;default:null;index" json:"locationId"`
	CapturedAt time.Time `gorm:"column:capturedAt;not null;index" json:"capturedAt"`
	CreatedAt  time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	User       User      `gorm:"foreignKey:UserID;references:ID;default:null" json:"user,omitempty"`
	Location   Location  `gorm:"foreignKey:LocationID;references:ID;default:null" json:"location,omitempty"`
}
