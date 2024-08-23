package subscribers

import (
	"log"

	"github.com/Tibz-Dankan/keep-active/internal/events"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/routes/request"
)

// Subscribes/listens to the topic "makeRequest"
// and performs an http get request to each app
// whose data has been received.
func subscribeToRequestEvent() {
	appCh := make(chan events.DataEvent)
	events.EB.Subscribe("makeRequest", appCh)
	type App = models.App

	for {
		appEvent := <-appCh
		app, ok := appEvent.Data.(App)

		if !ok {
			log.Println("Interface does not hold type App")
			return
		}

		go request.MakeAppRequest(app)
	}
}
