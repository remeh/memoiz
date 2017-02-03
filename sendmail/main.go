package main

import (
	"time"

	"remy.io/memoiz/config"
	"remy.io/memoiz/log"
	"remy.io/memoiz/storage"
)

const (
	CategoryReminderEmail = "CategoryReminderEmail"
	CategoryEnrichedEmail = "CategoryEnrichedEmail"

	RunFrequency     = time.Minute
	EmailFrequency   = time.Hour * 24
	EmailFrequencyPg = "1 day"
	EmailFirstAfter  = "1 day"
	//RunFrequency     = time.Second * 10
	//EmailFrequency   = time.Minute * 3
	//EmailFrequencyPg = "3 minute"
	//EmailFirstAfter  = "3 minute"
)

func main() {

	ticker := time.NewTicker(RunFrequency)

	log.Info("sendmail: starting")

	if err := prepare(); err != nil {
		log.Error("sendmail:", err)
	}

	if err := enrichEmailing(); err != nil {
		log.Error(err)
	}

	for t := range ticker.C {
		log.Debug("sendmail: waking up", t)
		memos, err := fetch()
		if err != nil {
			log.Error("sendmail:", err)
		}

		if len(memos) == 0 {
			continue
		}

		if err := send(memos, t); err != nil {
			log.Error("sendmail", err)
		}
	}
}

func prepare() error {
	_, err := storage.Init(config.Config.ConnString)
	return err
}
