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

func enrichEmailing(t time.Time) error {
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
		// memos for which we've find enrich infos
		// and that we'll send to the user.
		var toSend memos.Memos
		// found information
		var results mind.EnrichResults
		// memos for which we've looked for enrich infos
		// but we didn't find anything: we still want to
		// store the information that we've tried to send them.
		var lookedButNotFound memos.Memos

		// retrieve memos to do for this user
		// ----------------------

		if ms, err = enrichableMemos(uid, IntervalBetweenEachSend); err != nil {
			return log.Err("enrichEmailing", err)
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
			} else {
				lookedButNotFound = append(lookedButNotFound, m)
			}

			if len(toSend) >= MaxEnrichedPerMail {
				break
			}
		}

		// no memos to send to this user
		if len(toSend) == 0 {
			continue
		}

		log.Info("Sending Enriched Email to", uid)

		// gets this user account
		// ----------------------
		var acc accounts.SimpleUser

		if acc, _, err = accounts.DAO().UserByUid(uid); err != nil {
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

		allLooked := append(toSend, lookedButNotFound...)

		if err := memos.DAO().UpdateLastEmail(uid, allLooked.Uids(), t); err != nil {
			log.Error("enrichEmailing:", err)
			continue
		}

		if err := emailSent(acc, CategoryEnrichedEmail, t); err != nil {
			return log.Err("enrichEmailing", err)
		}

		// send the mail
		// ----------------------

		if err := email.SendEnrichedMemos(acc, toSend, results); err != nil {
			return log.Err("enrichEmailing", err)
		}
	}

	return nil
}
