package main

import (
	"fmt"
	"time"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/log"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/notify/email"
	"remy.io/memoiz/uuid"
)

// fetch fetches Ids of memos for which notification
// has never been done.
func fetch() (map[string]memos.Memos, error) {

	// first we want to retrieve for whom we'll
	// send some emails
	// ----------------------

	var err error
	var uids uuid.UUIDs

	if uids, err = getOwners(CategoryReminderEmail, EmailFrequency, 5); err != nil {
		return nil, log.Err("fetch", err)
	}

	// gets the memos of these owners
	// ----------------------

	if len(uids) == 0 {
		return make(map[string]memos.Memos), nil
	}

	return getRecentMemos(uids)
}

// send sends, per user, a list of memos per email
// in order to remind the user they've added them.
func send(memos map[string]memos.Memos, t time.Time) error {
	for owner, memos := range memos {
		memos = memos
		log.Info("Sending for", owner)

		// get the user
		// ----------------------

		var uid uuid.UUID
		var err error

		if uid, err = uuid.Parse(owner); err != nil {
			return log.Err("send", err)
		}

		var acc accounts.SimpleUser

		if acc, _, err = accounts.DAO().UserByUid(uid); err != nil {
			return fmt.Errorf("send: unknown user %q", owner)
		}

		// send the email
		// ----------------------

		if err := email.SendCategoryMail(acc, memos.GroupByCategory()); err != nil {
			return log.Err("send", err)
		}

		// store that the email has been sent.
		// ----------------------

		if err := emailSent(acc, CategoryReminderEmail, t); err != nil {
			return log.Err("send", err)
		}
	}

	return nil
}
