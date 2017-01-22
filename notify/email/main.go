package main

import (
	"fmt"
	"time"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/cards"
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
		cards, err := fetch()
		if err != nil {
			log.Error("notify/email:", err)
		}

		if len(cards) == 0 {
			continue
		}

		if err := send(cards, t); err != nil {
			log.Error("notify/email", err)
		}
	}
}

func prepare() error {
	_, err := storage.Init(config.Config.ConnString)
	return err
}

// fetch fetches Ids of cards for which notification
// has never been done.
func fetch() (map[string]cards.Cards, error) {

	// first we want to retrieve for whom we'll
	// send some emails
	// ----------------------

	var err error
	var uids uuid.UUIDs

	if uids, err = getOwners(EmailFrequency, 5); err != nil {
		return nil, log.Err("fetch", err)
	}

	// gets the cards of these owners
	// ----------------------

	if len(uids) == 0 {
		return make(map[string]cards.Cards), nil
	}

	return getCards(uids)
}

// TODO(remy): comment me.
func send(cards map[string]cards.Cards, t time.Time) error {
	for owner, cards := range cards {
		cards = cards
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

		if err := notify.SendCategoryMail(acc, cards.GroupByCategory()); err != nil {
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
