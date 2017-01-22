package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"

	"remy.io/scratche/accounts"
	"remy.io/scratche/cards"
	"remy.io/scratche/log"
	"remy.io/scratche/mind"
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
		SELECT u."uid", coalesce(max(es."creation_time"), '1970-01-01')
		FROM "user" u
		LEFT JOIN "emailing_sent" es ON
			u."uid" = es."owner_uid"
		LEFT JOIN "emailing_unsubscribe" eu ON
			eu."owner_uid" = u."uid"
		WHERE
			-- send the first mail 1 day after the
			-- user subscription
			u."creation_time" + interval '`+EmailFirstAfter+`' < now()
			AND
			-- send only to people not unsubscribed
			eu."creation_time" IS NULL
		GROUP BY u."uid"
		ORDER BY max(es."creation_time") DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, log.Err("getOwners", err)
	}

	uids := make(uuid.UUIDs, 0)

	defer rows.Close()
	for rows.Next() {
		var uid uuid.UUID
		var nt pq.NullTime

		if err := rows.Scan(&uid, &nt); err != nil {
			return nil, log.Err("getOwners: Scan", err)
		}

		var t time.Time
		if _, err := nt.Value(); nt.Valid && err == nil {
			t = nt.Time
		}

		if time.Since(t) > d {
			uids = append(uids, uid)
		}
	}

	return uids, nil
}

// getCards returns the cards per owners.
func getCards(owners uuid.UUIDs) (map[string]cards.Cards, error) {
	var rows *sql.Rows
	var err error

	if len(owners) == 0 {
		return nil, fmt.Errorf("notify/email: getCards: called with len(owners) == 0")
	}

	// query
	// ----------------------

	// parameters

	p := make([]interface{}, len(owners))
	for i, uid := range owners {
		p[i] = uid
	}

	// build in clause

	in := "("
	for i := range p {
		in += fmt.Sprintf("$%d", i+1)
		switch {
		case i < len(p)-1:
			in += ","
		case i == len(p)-1:
			in += ")"
		}
	}

	// finally query cards created between last mail and this mail.

	if rows, err = storage.DB().Query(`
		SELECT "owner_uid", array_agg("uid"), array_agg(text), array_agg("r_category")
		FROM "card"
		WHERE
			"owner_uid" IN `+in+`
			AND
			-- cards created between last mail and this email
			"creation_time" + interval '`+EmailFrequencyPg+`' > now()
		GROUP BY "owner_uid"
	`, p...); err != nil {
		return nil, log.Err("getCards", err)
	}

	if rows == nil {
		return make(map[string]cards.Cards), nil
	}

	// read the results
	// ----------------------

	rv := make(map[string]cards.Cards)

	defer rows.Close()
	for rows.Next() {
		var uid string
		var uids uuid.UUIDs
		var texts []string
		var cats []int64

		if err := rows.Scan(&uid, pq.Array(&uids), pq.Array(&texts), pq.Array(&cats)); err != nil {
			log.Error("notify/email: getCards:", err, "Continuing.")
			continue
		}

		if len(uids) != len(cats) || len(uids) != len(texts) {
			log.Error("notify/email: getCards: len(uids) != len(cats) for", uid, "Continuing.")
			continue
		}

		cards := make(cards.Cards, len(uids))
		for i, uid := range uids {
			cards[i].Uid = uid
			cards[i].CardRichInfo.Category = mind.Category(cats[i])
			if len(texts[i]) > 140 {
				cards[i].Text = texts[i][:140] + "..."
			} else {
				cards[i].Text = texts[i]
			}
		}

		rv[uid] = cards
	}

	return rv, nil
}

// emailSent stores in the database that an email has been sent
// to the given user at the given time.
func emailSent(acc accounts.SimpleUser, t time.Time) error {
	if _, err := storage.DB().Exec(`
		INSERT INTO "emailing_sent"
		("uid", "owner_uid", "type", "creation_time")
		VALUES
		($1, $2, $3, $4)
	`, uuid.New(), acc.Uid, CategoryReminderEmail, t); err != nil {
		return log.Err("emailSent", err)
	}
	return nil
}
