package events

import "log"

func InitEventSubscribers() {
	log.Println("Initiating event subscribers...")
	subscribeToPermissions()
}
