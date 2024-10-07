package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (f *BugReport) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (f *BugReport) Create(bugReport BugReport) (BugReport, error) {
	result := db.Create(&bugReport)

	if result.Error != nil {
		return bugReport, result.Error
	}

	bugReport, err := f.FindOne(bugReport.ID)
	if err != nil {
		return bugReport, err
	}

	return bugReport, nil
}

func (f *BugReport) Update() error {
	db.Save(&f)

	return nil
}

// Gets bug report by id
func (f *BugReport) FindOne(id string) (BugReport, error) {
	var bugReport BugReport

	db.First(&bugReport, "id = ?", id)

	return bugReport, nil
}

// Gets paginated bug reports for a specific user
func (f *BugReport) FindByUser(userId string, createdAtBeforeCursor time.Time) ([]BugReport, int64, error) {
	var bugReport []BugReport

	query := db.Table("feedbacks").Where("\"userId\" = ?", userId)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return bugReport, 0, err
	}

	if !createdAtBeforeCursor.IsZero() {
		query = query.Where("\"createdAt\" < ?", createdAtBeforeCursor)
	}

	if err := query.Order("\"createdAt\" desc").Limit(10).Find(&bugReport).Error; err != nil {
		return bugReport, 0, err
	}

	return bugReport, count, nil
}

// Gets all bug reports for a specific user
func (f *BugReport) FindAllByUser(userId string) ([]BugReport, error) {
	var bugReport []BugReport

	if err := db.Order("\"createdAt\" desc").Find(&bugReport, "\"userId\" = ?", userId).Error; err != nil {
		return bugReport, err
	}

	return bugReport, nil
}

// Gets paginated bug reports for all users
func (f *BugReport) FindAll(createdAtBeforeCursor time.Time) ([]BugReport, int64, error) {
	var bugReports, userBugReport []BugReport

	query := db.Table("bug_reports")

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if !createdAtBeforeCursor.IsZero() {
		query = query.Where("\"createdAt\" < ?", createdAtBeforeCursor)
	}

	if err := query.Order("\"createdAt\" desc").Limit(10).Find(&bugReports).Error; err != nil {
		return nil, 0, err
	}

	for _, bugReport := range bugReports {
		var user User
		db.Select("id", "name", "email", "role", "\"createdAt\"", "\"updatedAt\"").First(&user, "id = ?", bugReport.UserID)
		bugReport.User = user

		userBugReport = append(userBugReport, bugReport)
	}

	return userBugReport, count, nil
}
