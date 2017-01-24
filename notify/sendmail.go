package notify

import (
	"bytes"
	"fmt"
	"net/smtp"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/config"
	"remy.io/memoiz/log"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/mind"
	"remy.io/memoiz/notify/template"
)

var (
	UseMail = false

	Sender = "memos@mail.memoiz.com"
)

func init() {
	if len(config.Config.SmtpHost) != 0 &&
		config.Config.SmtpPort != 0 &&
		len(config.Config.SmtpLogin) != 0 &&
		len(config.Config.SmtpPassword) != 0 {

		if template.Root != nil {
			UseMail = true
			log.Info("Mailing activated.")
		}
	}
}

// ----------------------

type scmParam struct {
	SimpleUser accounts.SimpleUser
	Memos      map[mind.Category]memos.Memos
}

// SendCategoryMail sends an email to the given email
// to remind him he has recently added some new memos.
func SendCategoryMail(acc accounts.SimpleUser, cs map[mind.Category]memos.Memos) error {
	if !UseMail {
		return nil
	}

	host := fmt.Sprintf("%s:%d", config.Config.SmtpHost, config.Config.SmtpPort)

	auth := smtp.PlainAuth("",
		config.Config.SmtpLogin,
		config.Config.SmtpPassword,
		config.Config.SmtpHost,
	)

	buff := bytes.Buffer{}

	buff.WriteString("MIME-version: 1.0;\r\n")
	buff.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n")
	buff.WriteString(fmt.Sprintf("To: %s\r\n", acc.Email)) // TODO(remy): put the real email here.
	buff.WriteString("Subject: Hello!\r\n\r\n")            // TODO(remy): generate a subject

	html := template.Root.Lookup("base.html")
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

	err := smtp.SendMail(host, auth, Sender, []string{acc.Email}, buff.Bytes())
	if err != nil {
		return log.Err("SendCategoryMail", err)
	}

	return nil
}
