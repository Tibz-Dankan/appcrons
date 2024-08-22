package subscribers

import "log"

func InitEventSubscribers() {
	log.Println("Initiating global event subscribers...")
	go subscribeToRequestEvent()
	subscribeToPermissions()
}
