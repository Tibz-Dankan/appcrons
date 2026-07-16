package models

import "log"

// runUnknownUserBackfill ensures the shared placeholder "unknown user" row
// exists, then repoints every Location/SiteVisit row whose userId was left
// blank by an anonymous request (see middlewares.OptionalAuth) at that
// placeholder's id, and clears any SiteVisit.locationId left blank by a
// failed geo-IP lookup. It must run strictly between the "AutoMigrate with
// constraints disabled" and "AutoMigrate with constraints re-enabled"
// passes in Db(), so the userId/locationId columns are clean by the time
// GORM tries to create their foreign key constraints.
func runUnknownUserBackfill() error {
	user := User{}
	unknownUser, err := user.FindOrCreateUnknown()
	if err != nil {
		return err
	}
	log.Println("Unknown user ready, id:", unknownUser.ID)

	location := Location{}
	locationRows, err := location.BackfillMissingUserID(unknownUser.ID)
	if err != nil {
		return err
	}
	log.Println("Backfilled Location.userId rows:", locationRows)

	siteVisit := SiteVisit{}
	siteVisitUserRows, err := siteVisit.BackfillMissingUserID(unknownUser.ID)
	if err != nil {
		return err
	}
	log.Println("Backfilled SiteVisit.userId rows:", siteVisitUserRows)

	siteVisitLocationRows, err := siteVisit.BackfillMissingLocationID()
	if err != nil {
		return err
	}
	log.Println("Backfilled SiteVisit.locationId rows:", siteVisitLocationRows)

	return nil
}
