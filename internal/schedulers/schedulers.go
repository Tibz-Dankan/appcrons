package schedulers

import "log"

func InitSchedulers() {
	log.Println("Initiating schedulers...")
	go schedulePublishRequest()
	cleanUserAppMemory()
}
