package email

import (
	"bytes"
	"fmt"
	"net/smtp"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/log"
	"remy.io/memoiz/notify/template"
)

type prmParam struct {
	Token      string
	SimpleUser accounts.SimpleUser
}

// SendPasswordResetMail sends an email to the given email
// to reset a password.
func SendPasswordResetMail(su accounts.SimpleUser, tok string) error {
	if !UseMail {
		return nil
	}

	buff := bytes.Buffer{}

	// headers
	mailHeader(&buff, su.Email, "Reset your password")

	// content
	html := template.Root.Lookup("password_reset_mail.html")
	if html == nil {
		return fmt.Errorf("SendPasswordResetMail: can't find base template")
	}

	p := prmParam{
		SimpleUser: su,
		Token:      tok,
	}

	if err := html.Execute(&buff, p); err != nil {
		return log.Err("SendPasswordResetMail", err)
	}

	buff.WriteString("\r\n")

	// send
	err := smtp.SendMail(host(), auth(), Sender, []string{su.Email}, buff.Bytes())
	if err != nil {
		return log.Err("SendPasswordResetMail", err)
	}

	return nil
}
