package models

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/constants"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := u.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (u *User) Create(user User) (string, error) {
	result := db.Create(&user)

	if result.Error != nil {
		return "", result.Error
	}
	return user.ID, nil
}

func (u *User) FindOne(id string) (User, error) {
	var user User
	db.First(&user, "id = ?", id)

	return user, nil
}

func (u *User) FindByEmail(email string) (User, error) {
	var user User
	db.First(&user, "email = ?", email)

	return user, nil
}

// FindUnknown looks up the shared placeholder "unknown user" row by its
// well-known sentinel email. Like FindByEmail, a not-found result is not
// treated as an error - callers must check the returned user's ID.
func (u *User) FindUnknown() (User, error) {
	return u.FindByEmail(constants.UNKNOWN_USER_EMAIL)
}

// FindOrCreateUnknown returns the shared placeholder "unknown user" row,
// creating it if it doesn't exist yet. Anonymous requests to
// anonymous-tolerant routes (see middlewares.OptionalAuth) are attributed
// to this user's id instead of leaving userId blank. Idempotent - safe to
// call on every process start.
func (u *User) FindOrCreateUnknown() (User, error) {
	existing, err := u.FindUnknown()
	if err != nil {
		return existing, err
	}
	if existing.ID != "" {
		return existing, nil
	}

	// Role must be "user" (see ValidRole) - an invented role would fail
	// permission checks elsewhere. Password is a random, never-shared
	// string: this account is never signed into.
	unknownUser := User{
		Name:     "Unknown User",
		Email:    constants.UNKNOWN_USER_EMAIL,
		Password: uuid.New().String(),
		Role:     "user",
	}

	id, err := u.Create(unknownUser)
	if err != nil {
		return User{}, err
	}
	unknownUser.ID = id

	return unknownUser, nil
}

func (u *User) FindAll() ([]User, error) {
	var users []User
	db.Find(&users)

	return users, nil
}

type UserWithAppCount struct {
	User
	AppCount int64 `json:"appCount"`
}

func (u *User) FindAllAndIncludeAppCount(limit float64, cursor string) ([]UserWithAppCount, error) {
	var users []User
	var usersWithCount []UserWithAppCount

	query := db.Model(&User{}).Order("\"createdAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastUser User
		if err := db.Select("\"createdAt\"").Where("id = ?",
			cursor).First(&lastUser).Error; err != nil {
			return usersWithCount, err
		}
		query = query.Where("\"createdAt\" < ?", lastUser.CreatedAt)
	}

	query.Find(&users)

	for _, user := range users {
		var appCount int64
		countResult := db.Model(&App{}).Where("\"userId\" = ?", user.ID).Count(&appCount)
		if countResult.Error != nil {
			return nil, countResult.Error
		}

		user.Password = ""
		user.PasswordResetToken = ""
		user.PasswordResetExpiresAt = time.Time{}

		userWithCount := UserWithAppCount{
			User:     user,
			AppCount: appCount,
		}
		usersWithCount = append(usersWithCount, userWithCount)
	}

	return usersWithCount, nil
}

func (u *User) FindCount() (int64, error) {
	startTime := time.Now()
	var userCount int64

	if err := db.Model(&User{}).Count(&userCount).Error; err != nil {
		return userCount, err
	}
	log.Println("Total user count:", userCount)
	log.Println("queryTimeMS:", int(time.Since(startTime).Milliseconds()))

	return userCount, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *User) Update() error {
	db.Save(&u)

	return nil
}

func (u *User) Delete(id string) error {
	db.Delete(&User{}, id)

	return nil
}

// ResetPassword is the method updates user's password in db
func (u *User) ResetPassword(password string) error {
	hashedPassword, err := u.HashPassword(password)
	if err != nil {
		return err
	}

	u.Password = hashedPassword
	u.PasswordResetExpiresAt = time.Now()
	db.Save(&u)

	return nil
}

func (u *User) PasswordMatches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// Converts plain text password into hashed string
func (u *User) HashPassword(plainTextPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (u *User) ValidRole(role string) bool {
	roles := []string{"user", "sys_admin"}

	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

func (u *User) SetRole(role string) error {
	isValidRole := u.ValidRole(role)

	if !isValidRole {
		return errors.New("invalid user role")
	}

	u.Role = role
	return nil
}

func (u *User) FindByPasswordResetToken(resetToken string) (User, error) {
	var user User
	hashedToken := sha256.New()
	hashedToken.Write([]byte(resetToken))
	hashedTokenByteSlice := hashedToken.Sum(nil)
	hashedTokenString := hex.EncodeToString(hashedTokenByteSlice)

	result := db.Where("\"passwordResetToken\" = ? AND \"passwordResetExpiresAt\" > ?", hashedTokenString, time.Now()).Find(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (u *User) CreatePasswordResetToken() (string, error) {
	resetToken := uuid.NewString()

	hashedToken := sha256.New()
	hashedToken.Write([]byte(resetToken))
	hashedTokenByteSlice := hashedToken.Sum(nil)
	hashedTokenString := hex.EncodeToString(hashedTokenByteSlice)

	u.PasswordResetToken = hashedTokenString
	u.PasswordResetExpiresAt = time.Now().Add(20 * time.Minute)

	db.Save(&u)

	return resetToken, nil
}
