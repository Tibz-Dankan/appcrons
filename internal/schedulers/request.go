package schedulers

import (
	"log"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/events/publishers"
	"github.com/Tibz-Dankan/keep-active/internal/services"
)

// Runs the PublishRequestEvent fn at
// start of every minute that is a multiple of 10
func schedulePublishRequest() {
	for {
		now := time.Now()
		nextMinute := now.Truncate(time.Minute).Add(time.Minute)
		sleepDuration := nextMinute.Sub(now)
		minute := now.Minute()

		if minute%5 == 0 && now.Second() == 0 {
			if services.IsMaintenanceActive() {
				log.Println("Skipping PublishRequestEvent: maintenance window active")
			} else {
				publishers.PublishRequestEvent()
			}
		}

		time.Sleep(sleepDuration)
	}
}
