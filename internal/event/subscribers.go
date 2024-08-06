package event

import "log"

func EventSubscribers() {
	log.Println("Subscribing to app events...")
	subscribeToPermissions()
	subscribeToUpdateApp()
	subscribeToCreateApp()
	requestEventSubscriber()
}
