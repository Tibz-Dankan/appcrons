package models

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

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
		
		userWithCount := UserWithAppCount{
			User:     user,
			AppCount: appCount,
		}
		usersWithCount = append(usersWithCount, userWithCount)
	}
	
	return usersWithCount, nil
}


// func (a *Article) FindAllByPostedAt(limit int, articleIDCursor string,
// 	dateCursor time.Time, offset int) ([]Article, int64, error) {
// 	var articles []Article
// 	var count int64
// 	query := db.Model(&Article{}).
// 		Preload("Author").
// 		Order("\"postedAt\" DESC").
// 		Limit(int(limit))

// 	if offset != 0 {
// 		query = query.Offset(offset)
// 	}

	// if articleIDCursor != "" {
	// 	var lastArticle Article
	// 	if err := db.Select("\"postedAt\"").Where("id = ?",
	// 		articleIDCursor).First(&lastArticle).Error; err != nil {
	// 		return nil, 0, err
	// 	}
	// 	query = query.Where("\"postedAt\" < ?", lastArticle.PostedAt)
	// }

	// if !dateCursor.IsZero() {
	// 	query = query.Where("\"postedAt\" <= ?", dateCursor)
	// }

	// if err := query.Count(&count).Error; err != nil {
	// 	return nil, 0, err
	// }

	// query.Find(&articles)

// 	return articles, count, nil
// }

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
