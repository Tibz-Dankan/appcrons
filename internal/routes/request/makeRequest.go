package request

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/event"
	"github.com/Tibz-Dankan/keep-active/internal/models"

	"github.com/Tibz-Dankan/keep-active/internal/services"
)

func StartRequestScheduler() {
	go requestPublishScheduler()
	go requestEventSubscriber()
}

// Runs the requestPublisher fn at
// start of every minute
func requestPublishScheduler() {
	for {
		now := time.Now()
		nextMinute := now.Truncate(time.Minute).Add(time.Minute)
		sleepDuration := nextMinute.Sub(now)
		seconds := now.Second()
		if seconds == 0 {
			requestPublisher()
		}

		time.Sleep(sleepDuration)
	}
}

// Subscribes/listens to all published
// app request events
func requestEventSubscriber() {
	appCh := make(chan event.DataEvent)
	event.EB.Subscribe("makeRequest", appCh)
	type App = models.App

	for {
		var wg sync.WaitGroup
		appEvent := <-appCh
		app, ok := appEvent.Data.(App)

		if !ok {
			log.Println("Interface does not hold type App")
			return
		}
		wg.Add(1)
		go func(app models.App) {
			defer wg.Done()
			MakeAppRequest(app)
		}(app)
		wg.Wait()
	}
}

func requestPublisher() {
	app := models.App{}

	apps, err := app.FindAll()
	if err != nil {
		log.Println("Error fetching apps:", err)
		return
	}

	if len(apps) == 0 {
		return
	}

	publishRequests(apps)
}

// publishes all app request events
func publishRequests(apps []models.App) {
	var wg sync.WaitGroup
	for _, app := range apps {
		wg.Add(1)
		go func(app models.App) {
			defer wg.Done()
			event.EB.Publish("makeRequest", app)
		}(app)
	}
	wg.Wait()
}

// Makes request for the app
func MakeAppRequest(app models.App) {
	ok, err := validateApp(app)
	if err != nil {
		log.Println("Error validating the app: ", err)
	}
	if !ok {
		log.Println("Couldn't make request: ", app.Name)
		return
	}

	event.EB.Publish("app", app)

	response, err := services.MakeHTTPRequest(app.URL)
	if err != nil {
		log.Println("Request error:", err)
		return
	}

	request := models.Request{
		AppID:      app.ID,
		StatusCode: response.StatusCode,
		Duration:   response.RequestTimeMS,
		StartedAt:  response.StartedAt,
	}

	requestId, err := request.Create(request)
	if err != nil {
		log.Println("Error saving request:", err)
	}

	request.ID = requestId

	app.Request = []models.Request{request}

	event.EB.Publish("app", app)
	event.EB.Publish("updateApp", app)
}

// Validates the app's eligibility for making requests
func validateApp(app models.App) (bool, error) {
	if app.IsDisabled {
		return false, nil
	}

	hasLastRequest := len(app.Request) > 0

	// Check and validate requestInterval
	if len(app.RequestTime) == 0 {
		if !hasLastRequest {
			return true, nil
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
			return true, nil
		}
		return false, nil
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
		// log.Printf("app:%s rs:%s startTime:%s re:%s endTime: %s\n", app.Name, rt.Start, startTime, rt.End, endTime)

		isEqualToStartTime := currentTimeStart.Equal(startTime)
		isEqualToEndTime := currentTimeEnd.Equal(endTime)
		isGreaterThanStartTime := currentTimeStart.After(startTime)
		isLessThanEndTime := currentTimeEnd.Before(endTime)

		isWithinRequestTimeRange := isGreaterThanStartTime && isLessThanEndTime

		if isEqualToStartTime || isEqualToEndTime || isWithinRequestTimeRange {
			// Check and validate requestInterval
			log.Println("App time frame is correct")
			if !hasLastRequest {
				return true, nil
			}
			lastRequestStartedAt, _ := appDateStart.ISOTime()
			timeDiff := currentTimeStart.Sub(lastRequestStartedAt).Minutes()
			requestInterval, err := strconv.Atoi(app.RequestInterval)
			if err != nil {
				log.Println("Error converting string to integer:::", err)
			}

			if int(timeDiff) >= requestInterval {
				return true, nil
			}
		}
	}

	return false, nil
}
