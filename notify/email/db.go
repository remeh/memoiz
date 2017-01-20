package main

import (
	"time"

	"remy.io/scratche/log"
	"remy.io/scratche/storage"
	"remy.io/scratche/uuid"
)

// getOwners returns a given amount of owner
// which must receive a notification because last time they
// have been notified is bigger than the given duration.
func getOwners(d time.Duration, limit int) (uuid.UUIDs, error) {
	// query
	// ----------------------

	rows, err := storage.DB().Query(`
		SELECT u."uid", "emailing"."creation_time"
		FROM "user" u
		LEFT JOIN "emailing" e ON
			u."uid" = e."owner.uid"
		ORDER BY "emailing"."creation_time"
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, log.Err("getOwners", err)
	}

	uids := make(uuid.UUIDs, 0)

	defer rows.Close()
	for rows.Next() {
		var uid uuid.UUID
		var t time.Time

		if err := rows.Scan(&uid, &t); err != nil {
			return nil, log.Err("getOwners: Scan", err)
		}

		if time.Since(t) > d {
			uids = append(uids, uid)
		}
	}

	return uids, nil
}

// emailSent stores in the database that an email has been sent
// to the given user at the given time.
func emailSent(ownerUid uuid.UUID, t time.Time) error {
	if _, err := storage.DB().Exec(`
		INSERT INTO "emailing"
		("uid", "owner_uid", "type", "creation_time")
		VALUES
		($1, $2, $3, $4)
	`, uuid.New(), ownerUid, CategoryReminderEmail, t); err != nil {
		return log.Err("emailSent", err)
	}
	return nil
}
