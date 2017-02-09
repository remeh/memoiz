package email

import (
	"bytes"
	"fmt"
	"net/smtp"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/log"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/mind"
	"remy.io/memoiz/notify/template"
	"remy.io/memoiz/uuid"
)

type scmParam struct {
	SimpleUser accounts.SimpleUser
	Memos      map[mind.Category]memos.Memos
}

// SendCategoryMail sends an email to the given email
// to remind him he has recently added some new memos.
func SendRecentlyAddedMail(acc accounts.SimpleUser, cs map[mind.Category]memos.Memos, dumpDir string, sendUid uuid.UUID) error {
	if !UseMail {
		return nil
	}

	buff := bytes.Buffer{}

	// headers
	mailHeader(&buff, acc.Email, "Memos are waiting for you!")

	// content
	html := template.Root.Lookup("category_mail.html")
	if html == nil {
		return fmt.Errorf("SendCategoryMail: can't find base template")
	}

	p := scmParam{
		SimpleUser: acc,
		Memos:      cs,
	}

	if err := html.Execute(&buff, p); err != nil {
		log.Err("SendCategoryMail", err)
	}

	buff.WriteString("\r\n")

	dumpToFile(dumpDir, sendUid.String(), buff.Bytes())

	// send
	err := smtp.SendMail(host(), auth(), Sender, []string{acc.Email}, buff.Bytes())
	if err != nil {
		return log.Err("SendCategoryMail", err)
	}

	return nil
}
