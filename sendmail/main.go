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

	for t := range ticker.C {
		if err := categoryEmailing(t); err != nil {
			log.Error(err)
		}
		if err := enrichEmailing(t); err != nil {
			log.Error(err)
		}
	}
}

func prepare() error {
	_, err := storage.Init(config.Config.ConnString)
	return err
}
