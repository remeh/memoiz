package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/log"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/mind"
	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"
)

// getOwners returns a given amount of owner
// which must receive a notification because last time they
// have been notified with the given type of reminder
// is bigger than the given duration.
func getOwners(cat string, d time.Duration, limit int) (uuid.UUIDs, error) {
	// query
	// ----------------------

	// TODO(remy): send the email only to sub users.

	rows, err := storage.DB().Query(`
		SELECT u."uid", coalesce(max(es."creation_time"), '1970-01-01')
		FROM "user" u
		LEFT JOIN "emailing_sent" es ON
			u."uid" = es."owner_uid"
			AND
			es."type" = $1
		LEFT JOIN "emailing_unsubscribe" eu ON
			eu."owner_uid" = u."uid"
		WHERE
			-- send the first mail 1 day after the
			-- user subscription
			u."creation_time" + interval '`+FirstEmailAfter+`' < now()
			AND
			-- send only to people not unsubscribed
			eu."creation_time" IS NULL
		GROUP BY u."uid"
		ORDER BY max(es."creation_time")
		LIMIT $2
	`, cat, limit)
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

// enrichableMemos return enrichable memos available
// to be sent to their owner because they were not been
// sent since the given duration.
// Interval use the postgresql interval syntax.
func enrichableMemos(owner uuid.UUID, interval string) (memos.Memos, error) {
	var rows *sql.Rows
	var err error

	if owner == nil {
		return nil, fmt.Errorf("sendmail: enrichableMemos: nil owner provided")
	}

	if rows, err = storage.DB().Query(`
		SELECT "memo"."uid", text, "r_category", "r_title", "r_url"
		FROM "memo"
		LEFT JOIN "emailing_memo" em ON
			em."owner_uid" = "memo"."owner_uid"
			AND
			em."uid" = "memo"."uid"
		WHERE
			"memo"."owner_uid" = $1
			AND
			"memo"."state" = $2
			AND
			COALESCE(em."last_sent", "memo"."creation_time") + interval '`+interval+`'  < now()
		ORDER BY em."last_sent" DESC
	`, owner, memos.MemoActive); err != nil {
		return nil, log.Err("enrichableMemos", err)
	}

	if rows == nil {
		return make(memos.Memos, 0), nil
	}

	// read the results
	// ----------------------

	rv := make(memos.Memos, 0)

	defer rows.Close()
	for rows.Next() {
		var uid uuid.UUID
		var text, url, title string
		var cat int64

		if err := rows.Scan(&uid, &text, &cat, &title, &url); err != nil {
			log.Error("sendmail: enrichableMemos:", err, "Continuing.")
			continue
		}

		rv = append(rv, memos.Memo{
			Uid:  uid,
			Text: text,
			MemoRichInfo: memos.MemoRichInfo{
				Category: mind.Category(cat),
				Title:    title,
				Url:      url,
			},
		})
	}

	return rv, nil
}

// getReminderToSend returns memos which must be send
// due to the reminder.
func getReminderToSend(t time.Time, limit int) (memos.MemosMap, error) {
	var rows *sql.Rows
	var err error

	if rows, err = storage.DB().Query(`
		SELECT "memo"."owner_uid", array_agg("memo"."uid"), array_agg("text"), array_agg("r_category")
		FROM "memo"
		LEFT JOIN "emailing_memo" ON
			"emailing_memo"."uid" = "memo"."uid"
			AND
			"emailing_memo"."owner_uid" = "memo"."owner_uid"
		WHERE
			"reminder" IS NOT NULL
			AND
			COALESCE("last_sent", '1970-01-01 01:01'::timestamp with time zone) < "memo"."reminder"
			AND
			"reminder" < $1
		GROUP BY "memo"."owner_uid";
	`, t); err != nil {
		return nil, err
	}

	if rows == nil {
		return make(memos.MemosMap), nil
	}

	// read the results
	// ----------------------

	rv := make(memos.MemosMap)

	defer rows.Close()
	for rows.Next() {
		var uid string
		var uids uuid.UUIDs
		var texts []string
		var cats []int64

		if err := rows.Scan(&uid, pq.Array(&uids), pq.Array(&texts), pq.Array(&cats)); err != nil {
			log.Error("sendmail: getReminderToSend:", err, "Continuing.")
			continue
		}

		if len(uids) != len(cats) || len(uids) != len(texts) {
			log.Error("sendmail: getReminderToSend: len(uids) != len(cats) for", uid, "Continuing.")
			continue
		}

		memos := make(memos.Memos, len(uids))
		for i, uid := range uids {
			memos[i].Uid = uid
			memos[i].MemoRichInfo.Category = mind.Category(cats[i])
			memos[i].Text = texts[i]
		}

		rv[uid] = memos
	}

	return rv, nil
}

// getRecentMemos returns recent memos per owners
// recently created and not already sent to the owner.
func getRecentMemos(owners uuid.UUIDs) (memos.MemosMap, error) {
	var rows *sql.Rows
	var err error

	if len(owners) == 0 {
		return nil, fmt.Errorf("sendmail: getRecentMemos: called with len(owners) == 0")
	}

	// query
	// ----------------------

	// parameters

	p := make([]interface{}, len(owners))
	for i, uid := range owners {
		p[i] = uid
	}

	// build in clause

	in := "(" // TODO(remy): use InClause helper in storage pkg
	for i := range p {
		in += fmt.Sprintf("$%d", i+1)
		switch {
		case i < len(p)-1:
			in += ","
		case i == len(p)-1:
			in += ")"
		}
	}

	// finally query memos created between last mail and this mail.
	// TODO(remy): use a dynamic state instead of directly MemoActive

	if rows, err = storage.DB().Query(`
		SELECT "memo"."owner_uid", array_agg("memo"."uid"), array_agg(text), array_agg("r_category")
		FROM "memo"
		LEFT JOIN "emailing_memo" em ON
			em."owner_uid" = "memo"."owner_uid"
			AND
			em."uid" = "memo"."uid"
			AND
			-- we want all memos not send as "recently added"
			em."type" = '`+CategoryRecentlyAddedEmail+`'
		WHERE
			"memo"."owner_uid" IN `+in+`
			AND
			"state" = 'MemoActive'
			AND
			-- memos created between last email and this email
			(
					-- do not resend many times memos in "recently added"
					em."last_sent" IS NULL
					AND
					"memo"."creation_time" + interval '`+RecentlyAddedFrequencyPg+`' < now()
			)
		GROUP BY "memo"."owner_uid"
	`, p...); err != nil {
		return nil, log.Err("getRecentMemos", err)
	}

	if rows == nil {
		return make(memos.MemosMap), nil
	}

	// read the results
	// ----------------------

	rv := make(memos.MemosMap)

	defer rows.Close()
	for rows.Next() {
		var uid string
		var uids uuid.UUIDs
		var texts []string
		var cats []int64

		if err := rows.Scan(&uid, pq.Array(&uids), pq.Array(&texts), pq.Array(&cats)); err != nil {
			log.Error("sendmail: getRecentMemos:", err, "Continuing.")
			continue
		}

		if len(uids) != len(cats) || len(uids) != len(texts) {
			log.Error("sendmail: getRecentMemos: len(uids) != len(cats) for", uid, "Continuing.")
			continue
		}

		memos := make(memos.Memos, len(uids))
		for i, uid := range uids {
			memos[i].Uid = uid
			memos[i].MemoRichInfo.Category = mind.Category(cats[i])
			if len(texts[i]) > 140 {
				memos[i].Text = texts[i][:140] + "..."
			} else {
				memos[i].Text = texts[i]
			}
		}

		rv[uid] = memos
	}

	return rv, nil
}

// emailSent stores in the database that an email has been sent
// to the given user at the given time.
func emailSent(acc accounts.SimpleUser, sendUid uuid.UUID, cat string, t time.Time) error {
	if _, err := storage.DB().Exec(`
		INSERT INTO "emailing_sent"
		("uid", "owner_uid", "type", "creation_time")
		VALUES
		($1, $2, $3, $4)
	`, sendUid, acc.Uid, cat, t); err != nil {
		return log.Err("emailSent", err)
	}
	return nil
}
