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
	go RequestPublishScheduler()
	go requestEventSubscriber()
}

func RequestPublishScheduler() {
	for {
		RequestPublisher()
		time.Sleep(1 * time.Minute)
	}
}

func requestEventSubscriber() {
	appCh := make(chan event.DataEvent)
	event.EB.Subscribe("makeRequest", appCh)
	type App = models.App

	// Listening for events
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

func RequestPublisher() {
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

func validateApp(app models.App) (bool, error) {
	if app.IsDisabled {
		return false, nil
	}
	if app.RequestTime[0].ID == "" {
		// Check and validate requestInterval
		log.Println("App doesn't have requestTime")
		currentTime := time.Now()
		location := currentTime.Location().String()
		appDate := services.Date{TimeZone: location, ISOStringDate: app.Request[0].CreatedAt.String()}

		currentAppTime, _ := appDate.CurrentTime()
		lastRequestCreatedAt, _ := appDate.ISOTime()
		timeDiff := currentAppTime.Sub(lastRequestCreatedAt).Minutes()
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
		appDateStart := services.Date{TimeZone: rt.TimeZone, ISOStringDate: app.Request[0].CreatedAt.String(), HourMinSec: rt.Start}
		appDateEnd := services.Date{TimeZone: rt.TimeZone, ISOStringDate: app.Request[0].CreatedAt.String(), HourMinSec: rt.End}

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
			lastRequestCreatedAt, _ := appDateStart.ISOTime()
			timeDiff := currentTimeStart.Sub(lastRequestCreatedAt).Minutes()
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
