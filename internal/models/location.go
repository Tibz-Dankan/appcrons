package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/types"
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

// BackfillMissingUserID repoints every Location row whose userId was left
// blank ("" or NULL) at the given placeholder user id, so the
// userId->users(id) foreign key constraint can be created against clean
// data. Returns rows updated.
func (l *Location) BackfillMissingUserID(unknownUserId string) (int64, error) {
	result := db.Model(&Location{}).
		// Where("\"userId\" IS NULL OR \"userId\" = ''").
		Where("\"userId\" IS NULL").
		Update("userId", unknownUserId)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// FindByIP finds the most recently created Location for a given IP,
// used to avoid re-resolving geo-IP info for an IP already seen before.
func (l *Location) FindByIP(ip string) (Location, error) {
	var location Location

	db.Where("ip = ?", ip).Order("\"createdAt\" desc").First(&location)

	return location, nil
}

// // CountryDistribution is the per-country breakdown of unique users,
// // aggregated from geo-IP data captured in Location.Info.
// type CountryDistribution struct {
// 	Country     string `json:"country"`
// 	CountryCode string `json:"countryCode"`
// 	UserCount   int64  `json:"userCount"`
// }

// // FindUserCountryDistribution groups users by the country of their
// // earliest (signup-time) Location row. Info is unstructured JSON text
// // (no queryable country column), so rows are fetched ordered by
// // userId,createdAt asc and deduped/parsed in Go - this also keeps the
// // query portable across the Postgres/SQLite drivers this app runs on,
// // avoiding a Postgres-only "DISTINCT ON".
// func (l *Location) FindUserCountryDistribution() ([]CountryDistribution, error) {
// 	var locations []Location
// 	// err := db.Where("\"userId\" IS NOT NULL AND \"userId\" != ''").
// 	err := db.Where("\"userId\" IS NOT NULL").
// 		Order("\"userId\" asc").Order("\"createdAt\" asc").Find(&locations).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Println("locations: ", locations)

// 	type ipInfo struct {
// 		Country     string `json:"country"`
// 		CountryCode string `json:"countryCode"`
// 	}

// 	seenUsers := make(map[string]bool)
// 	counts := make(map[string]*CountryDistribution)

// 	for _, loc := range locations {
// 		if seenUsers[loc.UserID] {
// 			continue
// 		}
// 		seenUsers[loc.UserID] = true

// 		var info ipInfo
// 		_ = json.Unmarshal([]byte(loc.Info), &info)

// 		key := info.CountryCode
// 		country := info.Country
// 		if key == "" {
// 			key = "Unknown"
// 			country = "Unknown"
// 		}

// 		if _, ok := counts[key]; !ok {
// 			counts[key] = &CountryDistribution{Country: country, CountryCode: info.CountryCode}
// 		}
// 		counts[key].UserCount++
// 	}

// 	result := make([]CountryDistribution, 0, len(counts))
// 	for _, v := range counts {
// 		result = append(result, *v)
// 	}

// 	return result, nil
// }

type CountryDistribution struct {
	UserCount int64        `json:"userCount"`
	Users     []User       `json:"users"`
	LocInfo   types.IPInfo `json:"locInfo,omitempty"`
}

func (l *Location) FindUserCountryDistribution() ([]CountryDistribution, error) {
	startTime := time.Now()
	var locations []Location
	var distribution []CountryDistribution
	var user User
	var knownLocUserIDs []string
	// err := db.Where("\"userId\" IS NOT NULL").
	// 	Preload("User", func(db *gorm.DB) *gorm.DB {
	// 		return db.Select("id, name, email, role,\"createdAt\", \"updatedAt\"")
	// 	}).
	// 	Order("\"createdAt\" asc").
	// 	Distinct("\"userId\"").
	// 	Find(&locations).Error
	// if err != nil {
	// 	return nil, err
	// }
	err := db.Select(`DISTINCT ON ("userId") *`).
		Where(`"userId" IS NOT NULL`).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email, role,\"createdAt\", \"updatedAt\"")
		}).
		Order(`"userId"`).
		Order(`"createdAt" asc`).
		Find(&locations).Error
	if err != nil {
		return nil, err
	}

	log.Printf("locations: %+v\n", locations)

	seenCountries := make(map[string]bool)
	seenUsers := make(map[string]bool)

	for _, location := range locations {
		var info types.IPInfo
		if err := json.Unmarshal([]byte(location.Info), &info); err != nil {
			log.Println("Error unmarshalling IP info:", err)
			return nil, err
		}

		if seenCountries[info.CountryCode] {
			continue
		}
		seenCountries[info.CountryCode] = true

		// get all users from this country
		var currDistribution CountryDistribution
		var currLocUsers []User
		for _, location := range locations {
			if seenUsers[location.UserID] {
				continue
			}
			seenUsers[location.UserID] = true
			knownLocUserIDs = append(knownLocUserIDs, location.UserID)
			currLocUsers = append(currLocUsers, location.User)
		}

		currDistribution.LocInfo = info
		currDistribution.UserCount = int64(len(currLocUsers))
		currDistribution.Users = currLocUsers

		distribution = append(distribution, currDistribution)
	}

	// Get users from unkown location
	var unknownDistribution CountryDistribution
	unknownUsers, err := user.FindByExcludingIDs(knownLocUserIDs)
	if err != nil {
		return nil, err
	}
	unknownDistribution.UserCount = int64(len(unknownUsers))
	unknownDistribution.Users = unknownUsers
	unknownDistribution.LocInfo = types.IPInfo{
		Country:     "Unknown",
		CountryCode: "",
	}
	distribution = append(distribution, unknownDistribution)

	log.Println("queryTimeMS:", int(time.Since(startTime).Milliseconds()))

	return distribution, nil
}
