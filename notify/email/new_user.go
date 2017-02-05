package email

import (
	"bytes"
	"fmt"
	"net/smtp"

	"remy.io/memoiz/log"
	"remy.io/memoiz/notify/template"
)

// SendNewUserMail sends the email for a new subscription.
// TODO(remy): send the email validation email.
func SendNewUserMail(firstname, email string) error {
	sendMyselfNewUserMail(firstname, email)
	return sendUserSignupMail(firstname, email)
}

// ----------------------

func sendUserSignupMail(firstname, email string) error {
	if !UseMail {
		return nil
	}

	buff := bytes.Buffer{}

	// headers
	mailHeader(&buff, email, "Welcome to Memoiz")

	// content
	html := template.Root.Lookup("welcome_mail.html")
	if html == nil {
		return fmt.Errorf("sendUserSignupMail: can't find welcome template")
	}

	if err := html.Execute(&buff, nil); err != nil {
		return log.Err("sendUserSignupMail", err)
	}
	buff.WriteString("\r\n")

	// send
	err := smtp.SendMail(host(), auth(), Sender, []string{email}, buff.Bytes())
	if err != nil {
		return log.Err("sendUserSignupMail", err)
	}

	return nil
}

func sendMyselfNewUserMail(firstname, email string) error {
	if !UseMail {
		return nil
	}

	buff := bytes.Buffer{}

	// headers
	mailHeader(&buff, "me@remy.io", "New user!")

	// content
	buff.WriteString(fmt.Sprintf("<h1>New User</h1><ul><li>Email: %s</li><li>Firstname: %s</li></ul>\r\n", email, firstname))
	buff.WriteString("\r\n")

	// send
	err := smtp.SendMail(host(), auth(), Sender, []string{"me@remy.io"}, buff.Bytes())
	if err != nil {
		return log.Err("sendNewUserMail", err)
	}

	return nil
}
