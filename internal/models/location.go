package models

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (l *Location) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (l *Location) Create(location Location) (Location, error) {
	result := db.Create(&location)

	if result.Error != nil {
		return location, result.Error
	}

	return location, nil
}

func (l *Location) FindOne(id string) (Location, error) {
	var location Location

	db.First(&location, "id = ?", id)

	return location, nil
}

// FindByIP finds the most recently created Location for a given IP,
// used to avoid re-resolving geo-IP info for an IP already seen before.
func (l *Location) FindByIP(ip string) (Location, error) {
	var location Location

	db.Where("ip = ?", ip).Order("\"createdAt\" desc").First(&location)

	return location, nil
}

// CountryDistribution is the per-country breakdown of unique users,
// aggregated from geo-IP data captured in Location.Info.
type CountryDistribution struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	UserCount   int64  `json:"userCount"`
}

// FindUserCountryDistribution groups users by the country of their
// earliest (signup-time) Location row. Info is unstructured JSON text
// (no queryable country column), so rows are fetched ordered by
// userId,createdAt asc and deduped/parsed in Go - this also keeps the
// query portable across the Postgres/SQLite drivers this app runs on,
// avoiding a Postgres-only "DISTINCT ON".
func (l *Location) FindUserCountryDistribution() ([]CountryDistribution, error) {
	var locations []Location
	err := db.Where("\"userId\" IS NOT NULL AND \"userId\" != ''").
		Order("\"userId\" asc").Order("\"createdAt\" asc").Find(&locations).Error
	if err != nil {
		return nil, err
	}

	type ipInfo struct {
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
	}

	seenUsers := make(map[string]bool)
	counts := make(map[string]*CountryDistribution)

	for _, loc := range locations {
		if seenUsers[loc.UserID] {
			continue
		}
		seenUsers[loc.UserID] = true

		var info ipInfo
		_ = json.Unmarshal([]byte(loc.Info), &info)

		key := info.CountryCode
		country := info.Country
		if key == "" {
			key = "Unknown"
			country = "Unknown"
		}

		if _, ok := counts[key]; !ok {
			counts[key] = &CountryDistribution{Country: country, CountryCode: info.CountryCode}
		}
		counts[key].UserCount++
	}

	result := make([]CountryDistribution, 0, len(counts))
	for _, v := range counts {
		result = append(result, *v)
	}

	return result, nil
}
