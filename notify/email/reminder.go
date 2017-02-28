package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"sort"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/log"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/mind"
	"remy.io/memoiz/notify/template"
	"remy.io/memoiz/uuid"
)

// SendReminderMemos sends to the given user the list of memos
// as reminder.
// It also stores the email in the dumpDir directory using
// the sendUid as filename.
func SendReminderMemos(acc accounts.SimpleUser, ms memos.Memos, infos mind.EnrichResults, dumpDir string, sendUid uuid.UUID) error {
	if !UseMail {
		return nil
	}

	if len(ms) == 0 {
		return fmt.Errorf("SendReminderMemos: called with 0 memos")
	}

	// if we have 2 memos and only one of them has
	// an image, ensure it will be at the first position.
	// ----------------------

	em := buildEnrichedMemos(ms, infos)
	sort.Sort(ImageFirst(em))

	// build the email
	// ----------------------

	buff := bytes.Buffer{}

	// headers
	mailHeader(&buff, acc.Email, "Reminder: "+buildTitle(em))

	// content
	html := template.Root.Lookup("reminder_mail.html")
	if html == nil {
		return fmt.Errorf("SendReminderMemos: can't find reminder template")
	}

	p := semParam{
		SimpleUser: acc,
		Memos:      em,
	}

	if err := html.Execute(&buff, p); err != nil {
		return log.Err("SendReminderMemos", err)
	}

	buff.WriteString("\r\n")

	dumpToFile(dumpDir, sendUid.String(), buff.Bytes())

	// send
	err := smtp.SendMail(host(), auth(), Sender, []string{acc.Email}, buff.Bytes())
	if err != nil {
		return log.Err("SendReminderMemos", err)
	}

	return nil
}
