package event

import (
	"log"

	"github.com/Tibz-Dankan/keep-active/internal/models"
)

func requestEventSubscriber() {

	appCh := make(chan DataEvent)
	EB.Subscribe("makeRequest", appCh)
	type App = models.App

	// Listening for events
	for {
		appEvent := <-appCh
		app, ok := appEvent.Data.(App)

		if !ok {
			log.Println("Interface does not hold type App")
			return
		}

		go func(app models.App) {
			// request.MakeAppRequest(app)
		}(app)

		//TODO: make request here
	}
}
