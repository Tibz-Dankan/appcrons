package request

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/pubsub"
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
	fmt.Println("In the MakeRequest fn")

	apps, err := app.FindAll()
	if err != nil {
		fmt.Println("Error fetching apps:", err)
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
	psub := pubsub.PubSub{}

	if err := psub.Publish(app.ID, app); err != nil {
		log.Println("Error publishing:", err)
	}

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

	app.Request[0] = request

	if err := psub.Publish(app.ID, app); err != nil {
		log.Println("Error publishing:", err)
	}
}
