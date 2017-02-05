package main

import (
	"fmt"
	"os"
	"strings"
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

var (
	EmailDumpDir = ""
)

func main() {

	ticker := time.NewTicker(RunFrequency)

	log.Info("sendmail: starting")

	if err := prepare(); err != nil {
		log.Error("sendmail:", err)
		os.Exit(1)
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

	EmailDumpDir = os.Getenv("EMAIL_DUMP_DIR")
	if len(EmailDumpDir) == 0 {
		return log.Err("prepare:", fmt.Errorf("EMAIL_DUMP_DIR not set. Please set it."))
	}

	// remove trailing /
	if strings.HasSuffix(EmailDumpDir, "/") {
		EmailDumpDir = strings.TrimRight(EmailDumpDir, "/")
	}

	return err
}
