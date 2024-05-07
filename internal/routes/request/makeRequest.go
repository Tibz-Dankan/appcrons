package request

import (
	"log"
	"sync"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/event"
	"github.com/Tibz-Dankan/keep-active/internal/models"

	"github.com/Tibz-Dankan/keep-active/internal/services"
)

func StartRequestScheduler() {
	go MakeRequestScheduler()
}

func MakeRequestScheduler() {
	for {
		MakeRequest()
		time.Sleep(5 * time.Minute)
	}
}

func MakeRequest() {
	app := models.App{}

	apps, err := app.FindAll()
	if err != nil {
		log.Println("Error fetching apps:", err)
		return
	}

	if len(apps) == 0 {
		return
	}

	makeAllRequests(apps)
}

func makeAllRequests(apps []models.App) {
	var wg sync.WaitGroup

	for _, app := range apps {
		wg.Add(1)
		go func(app models.App) {
			defer wg.Done()
			makeRequest(app)
		}(app)
	}

	wg.Wait()
}

func makeRequest(app models.App) {
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
}
