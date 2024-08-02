package services

import (
	"log"
	"time"
)

// deletes all key-value pairs in the user app memory
// at 35 seconds for every minute that is a multiple of 5
func clearUserAppMemory() {
	log.Println("Inside clearUserAppMemory()")
	for {
		now := time.Now()
		minute := now.Minute()
		second := now.Second()

		var sleepDuration time.Duration
		if minute%5 == 0 && second <= 35 {
			sleepDuration = time.Duration(35-second) * time.Second
		} else {
			nextMinute := now.Truncate(time.Minute).Add(time.Minute)
			for nextMinute.Minute()%5 != 0 {
				nextMinute = nextMinute.Add(time.Minute)
			}
			sleepDuration = nextMinute.Sub(now) + (35 * time.Second)
		}

		time.Sleep(sleepDuration)

		now = time.Now()
		if now.Minute()%5 == 0 && now.Second() == 35 {
			UserAppMem.DeleteAll()
		}
	}
}

func StartClearUserAppMemoryScheduler() {
	go clearUserAppMemory()
}
