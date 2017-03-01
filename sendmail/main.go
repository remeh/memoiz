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
	CategoryRecentlyAddedEmail = "CategoryRecentlyAddedEmail"
	CategoryEnrichedEmail      = "CategoryEnrichedEmail"
	CategoryReminderEmail      = "CategoryReminderEmail"

	RunFrequency = time.Minute

	// RecentlyAddedFrequency is the frequency at which we send the recently
	// added emails.
	RecentlyAddedFrequency   = time.Hour * 24 * 3
	RecentlyAddedFrequencyPg = "3 day"

	// After how many time the very first email should be sent
	// to the user (not counting the subscription email).
	FirstEmailAfter = "1 day"
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
		if err := reminderEmailing(t); err != nil {
			log.Error(err)
		}
		if err := recentEmailing(t); err != nil {
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
