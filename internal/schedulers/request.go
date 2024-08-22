package schedulers

import (
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/events/publishers"
)

// Runs the PublishRequestEvent fn at
// start of every minute that is a multiple of 5
func schedulePublishRequest() {
	for {
		now := time.Now()
		nextMinute := now.Truncate(time.Minute).Add(time.Minute)
		sleepDuration := nextMinute.Sub(now)
		minute := now.Minute()

		if minute%5 == 0 && now.Second() == 0 {
			publishers.PublishRequestEvent()
		}

		time.Sleep(sleepDuration)
	}
}
