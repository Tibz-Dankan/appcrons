package publishers

import (
	"log"
	"sync"

	"github.com/Tibz-Dankan/keep-active/internal/events"
	"github.com/Tibz-Dankan/keep-active/internal/models"
)

// Publishes all apps to the topic "makeRequest"
// for the purposes of making requests to the apps.
// This function must be called in a scheduler
func PublishRequestEvent() {
	app := models.App{}

	apps, err := app.FindAll()
	if err != nil {
		log.Println("Error fetching apps:", err)
		return
	}

	if len(apps) == 0 {
		return
	}

	var wg sync.WaitGroup
	for _, app := range apps {
		wg.Add(1)
		go func(app models.App) {
			defer wg.Done()
			events.EB.Publish("makeRequest", app)
		}(app)
	}
	wg.Wait()
}
