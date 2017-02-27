// Email for programmed reminder.

package main

import (
	"fmt"
	"time"
)

// reminderEmailing fetches email to send because
// their reminder has been set and its time.
func reminderEmailing(t time.Time) error {

	// get memos for which the reminder has been
	// set and for which the last email sent in mode
	// 'reminder' is before this time.
	// ----------------------

	memos, err := getReminderToSend(t, 10)
	if err != nil {
		return err
	}

	fmt.Println(memos)

	return nil
}
