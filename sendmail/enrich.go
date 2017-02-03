package main

import (
	"time"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/log"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/mind"
	"remy.io/memoiz/notify/email"
	"remy.io/memoiz/uuid"
)

const (
	// Amount maximum of enriched memo
	// in an email.
	MaxEnrichedPerMail = 2
	// Frequency on which should be sent at maximum
	// each memo.
	// Meaning if this is '3 day' , should be sent
	// only each 3 day
	// Must use the postgresql interval syntax
	IntervalBetweenEachSend = "3 day"
)

func enrichEmailing() error {
	// TODO(remy): store last time when this user has received an email of this type
	// TODO(remy): store last time a memo has been sent to the user.
	// TODO(remy): store that this user has received a CategoryEnrichedEmail

	var uids uuid.UUIDs
	var err error

	if uids, err = getOwners(CategoryEnrichedEmail, time.Hour*24*2, 5); err != nil {
		return err
	}

	// for each users selected in the database,
	// try to find two memos for which we'll enrich the data
	for _, uid := range uids {
		var ms memos.Memos
		var err error
		var toSend memos.Memos
		var results mind.EnrichResults

		now := time.Now()

		// retrieve memos to do for this user
		// ----------------------

		if ms, err = enrichableMemos(uid, IntervalBetweenEachSend); err != nil {
			return err
		}

		// enrich these memos
		// ----------------------

		for _, m := range ms {
			var found bool
			var res mind.EnrichResult
			var err error

			if found, res, err = mind.Enrich(m.Text, m.Category); err != nil {
				return log.Err("enrichEmailing", err)
			}

			if found {
				toSend = append(toSend, m)
				results = append(results, res)
			}

			if len(toSend) >= MaxEnrichedPerMail {
				break
			}
		}

		// no memos to send to this user
		if len(toSend) == 0 {
			continue
		}

		// gets this user account
		// ----------------------
		var su accounts.SimpleUser

		if su, _, err = accounts.DAO().UserByUid(uid); err != nil {
			log.Error("enrichEmailing:", err)
			continue
		}

		// store that we have send this email
		// ----------------------
		//
		// NOTE(remy): we store it before actually sending it
		// because if one the update/insert fail, we will send zero time the email.
		// In the other order (update/insert after send), if the update fail,
		// we will send an infinite amount of time the email...

		if err := memos.DAO().UpdateLastEmail(uid, ms.Uids(), now); err != nil {
			log.Error("enrichEmailing:", err)
			continue
		}

		if err := emailSent(acc, CategoryReminderEmail, t); err != nil {
			return log.Err("send", err)
		}

		// send the mail
		// ----------------------

		if err := email.SendEnrichedMemos(su, toSend, results); err != nil {
			log.Error("enrichEmailing:", err)
		}
	}

	return nil
}
