package events

import (
	"log"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/models"
)

func subscribeToPermissions() {
	permissionCh := make(chan DataEvent)
	EB.Subscribe("permissions", permissionCh)

	permission := models.Permissions{}
	type User = models.User

	for {
		permissionEvent := <-permissionCh
		user, ok := permissionEvent.Data.(User)

		if !ok {
			log.Println("Couldn't find published permission data")
			return
		}

		startTime := time.Now()

		err := permission.Set(user.ID)
		if err != nil {
			log.Println("Error setting permissions:", err)
			return
		}

		log.Println("queryTimeMS:", int(time.Since(startTime).Milliseconds()))
		log.Println("Set user permissions successfully:", user.ID)
	}
}
