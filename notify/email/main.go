package main

import (
	"database/sql"
	"time"

	"remy.io/scratche/cards"
	"remy.io/scratche/config"
	"remy.io/scratche/log"
	"remy.io/scratche/mind"
	"remy.io/scratche/notify"
	"remy.io/scratche/storage"
	"remy.io/scratche/uuid"

	"github.com/lib/pq"
)

func main() {
	ticker := time.NewTicker(time.Second * 10)

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

		if err := analyze(cards); err != nil {
			log.Error("notify/email:", err)
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
	var rows *sql.Rows
	var err error

	// query
	// ----------------------

	if rows, err = storage.DB().Query(`
		SELECT "owner_uid", array_agg("uid"), array_agg(text), array_agg("r_category")
		FROM "card"
		WHERE "r_category" != 0
		GROUP BY "owner_uid"
	`); err != nil {
		return nil, log.Err("fetch", err)
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
			log.Error("notify/email: fetch:", err, "Continuing.")
			continue
		}

		if len(uids) != len(cats) || len(uids) != len(texts) {
			log.Error("notify/email: fetch: len(uids) != len(cats) for", uid, "Continuing.")
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

// analyze analyzer per owner a big set of cards.
func analyze(cards map[string]cards.Cards) error {
	for owner, cards := range cards {
		log.Info("Analyzing cards of", owner)

		notify.SendCategoryMail(cards.GroupByCategory())
	}

	return nil
}
