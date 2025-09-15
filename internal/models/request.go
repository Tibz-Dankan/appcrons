package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Request) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (r *Request) Create(request Request) (Request, error) {
	result := db.Create(&request)

	if result.Error != nil {
		return request, result.Error
	}
	return request, nil
}

func (r *Request) FindOne(id string) (Request, error) {
	var request Request

	if request.ID != "" {
		return request, nil
	}
	db.First(&request, "id = ?", id)

	return request, nil
}

func (r *Request) FindByApp(appId string, createdAtBeforeCursor time.Time) ([]Request, int64, error) {
	var requests []Request

	// 12-hour threshold
	threshold := time.Now().Add(-12 * time.Hour)

	query := db.Table("requests").Where("\"appId\" = ?", appId).Where("\"createdAt\" > ?", threshold)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if !createdAtBeforeCursor.IsZero() {
		query = query.Where("\"createdAt\" < ?", createdAtBeforeCursor)
	}

	if err := query.Order("\"createdAt\" desc").Limit(10).Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, count, nil
}

func (r *Request) FindAll() ([]Request, error) {
	var request []Request
	db.Find(&request)

	return request, nil
}

func (r *Request) FindCurrentMonthCount() (int64, error) {
	startTime := time.Now()
	var requestCount int64

	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstDayOfNextMonth := firstDayOfMonth.AddDate(0, 1, 0)

	if err := db.Model(&Request{}).
		Where("\"createdAt\" >= ? AND \"createdAt\" < ?", firstDayOfMonth, firstDayOfNextMonth).
		Count(&requestCount).Error; err != nil {
		return requestCount, err
	}

	currentMonth := now.Format("2006-01")
	log.Printf("Total Current '%s' count:%d", currentMonth, requestCount)
	log.Println("queryTimeMS:", int(time.Since(startTime).Milliseconds()))
	return requestCount, nil
}


func (r *Request) FindExistingCount() (int64, error) {
	startTime := time.Now()
	var existingRequestCount int64
	
	if err := db.Model(&RequestCount{}).
		Select("COALESCE(SUM(count), 0)").
		Scan(&existingRequestCount).Error; err != nil {
		return existingRequestCount, err
	}
	
	log.Printf("Total existing request count: %d", existingRequestCount)
	log.Println("queryTimeMS:", int(time.Since(startTime).Milliseconds()))
	return existingRequestCount, nil
}

func (r *Request) FindTotalCount() (int64, error) {
	var count int64

	currentMonthCount, err:= r.FindCurrentMonthCount()
	if  err != nil {
		return count, err
	}
	ExistingCount, err:= r.FindExistingCount()
	if  err != nil {
		return count, err
	}
	count = currentMonthCount+ExistingCount
	
	log.Printf("Final total request count: %d", count)
	return count, nil
}

func (r *Request) Delete(id string) error {

	if err := db.Unscoped().Where("id = ?", id).Delete(&Request{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Request) DeleteByApp(appId string) error {

	if err := db.Unscoped().Where("\"appId\" = ?", appId).Delete(&Request{}).Error; err != nil {
		return err
	}

	return nil
}
