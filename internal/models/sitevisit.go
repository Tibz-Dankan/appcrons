package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *SiteVisit) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (s *SiteVisit) Create(siteVisit SiteVisit) (SiteVisit, error) {
	result := db.Create(&siteVisit)

	if result.Error != nil {
		return siteVisit, result.Error
	}

	return siteVisit, nil
}

func (s *SiteVisit) FindOne(id string) (SiteVisit, error) {
	var siteVisit SiteVisit

	db.First(&siteVisit, "id = ?", id)

	return siteVisit, nil
}

// BackfillMissingUserID repoints every SiteVisit row whose userId was left
// blank ("" or NULL) at the given placeholder user id, so the
// userId->users(id) foreign key constraint can be created against clean
// data. Returns rows updated.
func (s *SiteVisit) BackfillMissingUserID(unknownUserId string) (int64, error) {
	result := db.Model(&SiteVisit{}).
		// Where("\"userId\" IS NULL OR \"userId\" = ''").
		Where("\"userId\" IS NULL").
		Update("userId", unknownUserId)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// BackfillMissingLocationID nulls out every SiteVisit row whose locationId
// was left blank (e.g. postSiteVisit saved location.ID from
// services.GetUserLocationByIP even when that call failed and returned a
// zero-value Location). Unlike userId there's no "unknown location"
// placeholder to point at - a genuine SQL NULL is correct here (Postgres
// exempts NULL foreign key values from constraint validation), and is what
// lets the locationId->locations(id) constraint be created. Returns rows
// updated.
func (s *SiteVisit) BackfillMissingLocationID() (int64, error) {
	result := db.Model(&SiteVisit{}).
		// Where("\"locationId\" IS NULL OR \"locationId\" = ''").
		Where("\"locationId\" IS NULL").
		Update("locationId", nil)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
