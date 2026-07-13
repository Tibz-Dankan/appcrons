package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
)

type IPInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
	Query       string  `json:"query"`
}

// GetIPInfo resolves geo-IP info for the given IP via the free ip-api.com
// service (no API key required).
func GetIPInfo(ip string) (*IPInfo, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ip-api request failed: %s", res.Status)
	}

	var info IPInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}

// GetUserLocationByIP resolves (and caches) a Location row for the given IP.
// userID may be "" for anonymous requests. In testing/staging environments
// the live ip-api.com call is skipped to keep tests fast and offline.
func GetUserLocationByIP(userID string, ip string) (models.Location, error) {
	location := models.Location{}

	if ip == "" {
		return location, fmt.Errorf("empty client IP")
	}

	existing, err := location.FindByIP(ip)
	if err != nil {
		return location, err
	}
	if existing.ID != "" {
		return existing, nil
	}

	var infoJSON []byte

	if os.Getenv("GO_ENV") == "testing" || os.Getenv("GO_ENV") == "staging" {
		infoJSON, _ = json.Marshal(IPInfo{Query: ip})
	} else {
		info, err := GetIPInfo(ip)
		if err != nil {
			log.Println("Error fetching IP info:", err)
			info = &IPInfo{Query: ip}
		}
		infoJSON, err = json.Marshal(info)
		if err != nil {
			log.Println("Error marshalling IP info:", err)
			infoJSON = []byte("{}")
		}
	}

	return location.Create(models.Location{
		UserID: userID,
		IP:     ip,
		Info:   string(infoJSON),
	})
}
