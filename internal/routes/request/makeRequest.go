package request

import (
	"log"
	"strconv"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/events"
	"github.com/Tibz-Dankan/keep-active/internal/models"

	"github.com/Tibz-Dankan/keep-active/internal/services"
)

// Makes request for the app
func MakeAppRequest(app models.App) {
	isValid, requiresExternalRequest, err := validateApp(app)
	if err != nil {
		log.Println("Error validating the app: ", err)
	}
	if !isValid {
		log.Println("Couldn't make request: ", app.Name)
		return
	}
	response := services.Response{}
	appRequestProgress := services.AppRequestProgress{App: app, InProgress: true}
	services.UserAppMem.Add(app.UserID, appRequestProgress)

	events.EB.Publish("appRequestProgress", appRequestProgress)

	if !requiresExternalRequest {
		response, err = services.MakeHTTPRequest(app.URL)
		if err != nil {
			log.Println("Request error:", err)
			return
		}
	}

	if requiresExternalRequest {
		response, err = services.MakeExternalHTTPRequest(app.URL)
		if err != nil {
			log.Println("Request error:", err)
			return
		}
	}

	request := models.Request{
		AppID:      app.ID,
		StatusCode: response.StatusCode,
		Duration:   response.RequestTimeMS,
		StartedAt:  response.StartedAt,
	}
	log.Printf("Request statusCode: %d Duration: %d URL: %s\n", request.StatusCode, request.Duration, app.URL)

	request, err = request.Create(request)
	if err != nil {
		log.Println("Error saving request:", err)
	}

	appRequestProgress.Request = []models.Request{request}
	appRequestProgress.InProgress = false
	services.UserAppMem.Add(app.UserID, appRequestProgress)

	events.EB.Publish("appRequestProgress", appRequestProgress)
}

// Validates the app's eligibility for making requests
// It returns isValid, requiresExternalRequest and error
func validateApp(app models.App) (bool, bool, error) {
	if app.IsDisabled {
		return false, false, nil
	}

	requiresExternal, err := appRequiresExternalRequest(app)
	if err != nil {
		log.Println("Error getting the external determination", err)
	}

	hasLastRequest := len(app.Request) > 0

	// Check and validate requestInterval
	if len(app.RequestTime) == 0 {
		if !hasLastRequest {
			return true, requiresExternal, nil
		}
		log.Println("App doesn't have requestTime")
		currentTime := time.Now()
		location := currentTime.Location().String()
		appDate := services.Date{TimeZone: location, ISOStringDate: app.Request[0].StartedAt.String()}

		currentAppTime, _ := appDate.CurrentTime()
		lastRequestStartedAt, _ := appDate.ISOTime()
		timeDiff := currentAppTime.Sub(lastRequestStartedAt).Minutes()
		requestInterval, err := strconv.Atoi(app.RequestInterval)
		if err != nil {
			log.Println("Error converting string to integer:", err)
		}

		if int(timeDiff) >= requestInterval {
			return true, requiresExternal, nil
		}
		return false, requiresExternal, nil
	}

	for _, rt := range app.RequestTime {
		// Check and validate requestTime slot

		lastReqStartedAtStr := time.Now().String()
		if hasLastRequest {
			lastReqStartedAtStr = app.Request[0].StartedAt.String()
		}

		appDateStart := services.Date{TimeZone: rt.TimeZone, ISOStringDate: lastReqStartedAtStr, HourMinSec: rt.Start}
		appDateEnd := services.Date{TimeZone: rt.TimeZone, ISOStringDate: lastReqStartedAtStr, HourMinSec: rt.End}

		startTime, _ := appDateStart.HourMinSecTime()
		endTime, _ := appDateEnd.HourMinSecTime()
		currentTimeStart, _ := appDateStart.CurrentTime()
		currentTimeEnd, _ := appDateEnd.CurrentTime()

		isEqualToStartTime := currentTimeStart.Equal(startTime)
		isEqualToEndTime := currentTimeEnd.Equal(endTime)
		isGreaterThanStartTime := currentTimeStart.After(startTime)
		isLessThanEndTime := currentTimeEnd.Before(endTime)

		isWithinRequestTimeRange := isGreaterThanStartTime && isLessThanEndTime

		if isEqualToStartTime || isEqualToEndTime || isWithinRequestTimeRange {
			// Check and validate requestInterval
			log.Println("App time frame is correct")
			if !hasLastRequest {
				return true, requiresExternal, nil
			}
			lastRequestStartedAt, _ := appDateStart.ISOTime()
			timeDiff := currentTimeStart.Sub(lastRequestStartedAt).Minutes()
			requestInterval, err := strconv.Atoi(app.RequestInterval)
			if err != nil {
				log.Println("Error converting string to integer:::", err)
			}

			if int(timeDiff) >= requestInterval {
				return true, requiresExternal, nil
			}
		}
	}

	return false, requiresExternal, nil
}

func appRequiresExternalRequest(app models.App) (bool, error) {
	hasLastRequest := len(app.Request) > 0
	hasRequestTime := len(app.RequestTime) > 0

	if !hasLastRequest {
		return true, nil
	}

	if !hasRequestTime {
		currentTime := time.Now()
		location := currentTime.Location().String()
		appDate := services.Date{TimeZone: location, ISOStringDate: app.Request[0].StartedAt.String()}

		currentAppTime, _ := appDate.CurrentTime()
		lastRequestStartedAt, _ := appDate.ISOTime()
		timeDiff := currentAppTime.Sub(lastRequestStartedAt).Minutes()

		if int(timeDiff) > 15 {
			return true, nil
		}

		if app.Request[0].StatusCode == 503 && app.Request[0].Duration < 15*1000 {
			return true, nil
		}
	}

	for _, rt := range app.RequestTime {
		lastReqStartedAtStr := time.Now().String()
		if hasLastRequest {
			lastReqStartedAtStr = app.Request[0].StartedAt.String()
		}

		appDate := services.Date{TimeZone: rt.TimeZone, ISOStringDate: lastReqStartedAtStr, HourMinSec: rt.Start}
		currentAppTime, _ := appDate.CurrentTime()
		lastRequestStartedAt, _ := appDate.ISOTime()
		timeDiff := currentAppTime.Sub(lastRequestStartedAt).Minutes()

		if int(timeDiff) > 15 {
			return true, nil
		}

		if app.Request[0].StatusCode == 503 && app.Request[0].Duration < 15*1000 {
			return true, nil
		}
	}

	return false, nil
}
