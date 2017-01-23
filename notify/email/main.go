package main

import (
	"fmt"
	"time"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/config"
	"remy.io/memoiz/log"
	"remy.io/memoiz/notify"
	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"
)

const (
	CategoryReminderEmail = "CategoryReminderEmail"
	RunFrequency          = time.Minute
	EmailFrequency        = time.Hour * 24
	EmailFrequencyPg      = "1 day"
	EmailFirstAfter       = "1 day"
	//RunFrequency     = time.Second * 10
	//EmailFrequency   = time.Minute * 3
	//EmailFrequencyPg = "3 minute"
	//EmailFirstAfter  = "3 minute"
)

func main() {
	ticker := time.NewTicker(RunFrequency)

	log.Info("notify/email: starting")

	if err := prepare(); err != nil {
		log.Error("notify/email:", err)
	}

	for t := range ticker.C {
		log.Debug("notify/email: waking up", t)
		memos, err := fetch()
		if err != nil {
			log.Error("notify/email:", err)
		}

		if len(memos) == 0 {
			continue
		}

		if err := send(memos, t); err != nil {
			log.Error("notify/email", err)
		}
	}
}

func prepare() error {
	_, err := storage.Init(config.Config.ConnString)
	return err
}

// fetch fetches Ids of memos for which notification
// has never been done.
func fetch() (map[string]memos.Memos, error) {

	// first we want to retrieve for whom we'll
	// send some emails
	// ----------------------

	var err error
	var uids uuid.UUIDs

	if uids, err = getOwners(EmailFrequency, 5); err != nil {
		return nil, log.Err("fetch", err)
	}

	// gets the memos of these owners
	// ----------------------

	if len(uids) == 0 {
		return make(map[string]memos.Memos), nil
	}

	return getMemos(uids)
}

// TODO(remy): comment me.
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

		if err := notify.SendCategoryMail(acc, memos.GroupByCategory()); err != nil {
			return log.Err("send", err)
		}

		// store that the email has been sent.
		// ----------------------

		if err := emailSent(acc, t); err != nil {
			return log.Err("send", err)
		}
	}

	return nil
}
