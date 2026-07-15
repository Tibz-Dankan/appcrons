package auth

import (
	"log"

	"github.com/Tibz-Dankan/keep-active/internal/models"
)

// CreateUnknownUser idempotently seeds the shared placeholder "unknown
// user" row (see models.User.FindOrCreateUnknown) that
// middlewares.OptionalAuth falls back to for anonymous callers. Not an
// HTTP route - a startup bootstrap hook called explicitly from
// cmd/main.go, mirroring this repo's existing bootstrap style
// (schedulers.InitSchedulers(), subscribers.InitEventSubscribers(), etc).
//
// By the time main() calls this, models.FinishMigration() has already
// ensured the unknown user exists as part of its own backfill (see
// internal/models/db.go), so this call is normally a no-op find. It
// exists as the explicit, visible "auth handlers" entry point for it.
func CreateUnknownUser() {
	user := models.User{}
	unknownUser, err := user.FindOrCreateUnknown()
	if err != nil {
		log.Println("Error ensuring unknown user exists:", err)
		return
	}
	log.Println("Unknown user ready, id:", unknownUser.ID)
}
