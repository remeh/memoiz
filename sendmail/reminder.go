// Email for programmed reminder.

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

// reminderEmailing fetches email to send because
// their reminder has been set and its time.
func reminderEmailing(t time.Time) error {
	log.Debug("reminderEmailing: waking up", t)

	// get memos for which the reminder has been
	// set and for which the last email sent in mode
	// 'reminder' is before this time.
	// ----------------------

	mmap, err := getReminderToSend(t, 10)
	if err != nil {
		return err
	}

	for id, ms := range mmap {
		uid, err := uuid.Parse(id)
		if err != nil {
			return err
		}

		// enrich results
		results := make(mind.EnrichResults, len(ms))

		// gets this user account
		// ----------------------
		var acc accounts.SimpleUser

		if acc, _, err = accounts.DAO().UserByUid(uid); err != nil {
			log.Error("enrichEmailing:", err)
			continue
		}

		// we want to try to enrich this memo before sending
		// the reminder.
		// ----------------------

		for i, m := range ms {
			var found bool
			var res mind.EnrichResult
			var err error

			if found, res, err = mind.Enrich(m.Text, m.Category); err != nil {
				return log.Err("reminderEmailing", err)
			}

			if found {
				results[i] = res
			}
		}

		// update that we've sent this email
		// ----------------------

		if err := memos.DAO().UpdateLastEmail(uid, ms.Uids(), CategoryReminderEmail, t); err != nil {
			return log.Err("reminderEmailing:", err)
		}

		sendUid := uuid.New()

		if err := emailSent(acc, sendUid, CategoryReminderEmail, t); err != nil {
			return log.Err("reminderEmailing", err)
		}

		// delete the reminder
		// ----------------------

		if err := memos.DAO().DeleteReminders(uid, ms.Uids()); err != nil {
			return log.Err("reminderEmailing:", err)
		}

		// send the mail
		// ----------------------

		if err := email.SendReminderMemos(acc, ms, results, EmailDumpDir, sendUid); err != nil {
			return log.Err("reminderEmailing", err)
		}
	}

	return nil
}
