// Mail for "Recently added" memos.

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

func recentEmailing(t time.Time) error {
	log.Debug("categoryEmailing: waking up", t)
	memos, err := fetchRecent()
	if err != nil {
		return err
	}

	if len(memos) == 0 {
		return nil
	}

	if err := sendRecent(memos, t); err != nil {
		return err
	}

	return nil
}

// fetch fetches Ids of memos for which notification
// has never been done and that have been added recently.
func fetchRecent() (map[string]memos.Memos, error) {

	// first we want to retrieve for whom we'll
	// send some emails
	// ----------------------

	var err error
	var uids uuid.UUIDs

	if uids, err = getOwners(CategoryRecentlyAddedEmail, RecentlyAddedFrequency, 5); err != nil {
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
// in order to remind the user they've recently added them.
func sendRecent(memos map[string]memos.Memos, t time.Time) error {
	for owner, memos := range memos {
		memos = memos
		log.Info("Sending Recently Added Email to", owner)

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

		// store that the email has been sent.
		// ----------------------

		sendUid := uuid.New()

		if err := emailSent(acc, sendUid, CategoryRecentlyAddedEmail, t); err != nil {
			return log.Err("send", err)
		}

		// send the email
		// ----------------------

		if err := email.SendRecentlyAddedMail(acc, memos.GroupByCategory(), EmailDumpDir, sendUid); err != nil {
			return log.Err("send", err)
		}
	}

	return nil
}
