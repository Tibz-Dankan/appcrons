package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (u *OTP) BeforeCreate(tx *gorm.DB) error {

	uuid := uuid.New().String()
	OTP := "should be Random six digit number"

	tx.Statement.SetColumn("ID", uuid)
	tx.Statement.SetColumn("OPT", OTP)
	return nil
}

func (u *OTP) Create(opt OTP) (string, error) {

	// Expire all other existing user opts
	result := db.Create(&opt)

	if result.Error != nil {
		return "", result.Error
	}
	return opt.ID, nil
}

func (u *OTP) FindOne(id string) (OTP, error) {
	var opt OTP
	db.First(&opt, "id = ?", id)

	return opt, nil
}

func (u *OTP) FindByUser(userId string) (OTP, error) {
	var otp OTP
	db.First(&otp, "userId = ?", userId)

	return otp, nil
}

func (u *OTP) ExpireUserOPT(userId string) error {

	//TODO: to expire opt here
	return nil
}
