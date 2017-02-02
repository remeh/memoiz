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
	MaxEnrichedPerMail = 2
)

func enrichEmailing() error {
	// TODO(remy): pick maximum 2 of its notes, not recently sent to him, for which we can find content
	// TODO(remy): build an email using all these information and the enriched template
	// TODO(remy): send it this email

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

		// retrieve memos to do for this user

		if ms, err = enrichableMemos(uid, time.Hour); err != nil {
			return err
		}

		// nrich these memos

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

			if len(toSend) > MaxEnrichedPerMail {
				break
			}
		}

		// we've no memos to send to this user
		if len(toSend) == 0 {
			continue
		}

		var su accounts.SimpleUser

		// gets this user account

		if su, _, err = accounts.DAO().UserByUid(uid); err != nil {
			log.Error("enrichEmailing:", err)
			continue
		}

		// send the mail

		if err := email.SendEnrichedMemos(su, toSend, results); err != nil {
			log.Error("enrichEmailing:", err)
		}
	}

	return nil
}
