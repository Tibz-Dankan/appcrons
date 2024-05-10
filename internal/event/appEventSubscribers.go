package event

import (
	"log"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
)

func subscribeToUpdateApp() {

	appCh := make(chan DataEvent)

	EB.Subscribe("updateApp", appCh)

	type App = models.App

	// Listening for events
	for {
		appEvent := <-appCh

		publishedApp, ok := appEvent.Data.(App)
		app := models.App{}

		startTime := time.Now()

		latestAppData, err := app.FindOne(publishedApp.ID)
		if err != nil {
			log.Println("Error fetching latest app data:", err)
		}

		duration := time.Since(startTime)
		queryTimeMS := int(duration.Milliseconds())

		log.Println("queryTimeMS:", queryTimeMS)

		if !ok {
			log.Println("Interface does not hold type App")
			return
		}

		appCache := models.AppCache{}
		if err := appCache.WriteOneToAll(latestAppData); err != nil {
			log.Println("Error writing to all apps:", err)
		}
		// if err := appCache.WriteOneToUser(latestAppData); err != nil {
		// 	log.Println("Error writing to all apps to user:", err)
		// }
	}
}

func subscribeToCreateApp() {

	appCh := make(chan DataEvent)

	EB.Subscribe("createApp", appCh)

	type App = models.App

	// Listening for events
	for {
		appEvent := <-appCh

		app, ok := appEvent.Data.(App)

		if !ok {
			log.Println("Interface does not hold type App")
			return
		}

		app.RequestTime = []models.RequestTime{}
		app.Request = []models.Request{}

		appCache := models.AppCache{}
		if err := appCache.WriteOneToAll(app); err != nil {
			log.Println("Error writing to all apps:", err)
		}
		// if err := appCache.WriteOneToUser(app); err != nil {
		// 	log.Println("Error writing to all apps to user:", err)
		// }
	}
}
