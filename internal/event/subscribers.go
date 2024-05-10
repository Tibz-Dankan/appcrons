package event

import "log"

func EventSubscribers() {
	log.Println("Subscribing to app events...")
	subscribeToUpdateApp()
	subscribeToCreateApp()
}
