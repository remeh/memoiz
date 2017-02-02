package email

import (
	"bytes"
	"fmt"
	"net/smtp"

	"remy.io/memoiz/config"
	"remy.io/memoiz/log"
)

// SendNewUserMail sends the email for a new subscription.
// TODO(remy): send the email validation email.
func SendNewUserMail(firstname, email string) error {
	return sendUserSignupMail(firstname, email)
}

// ----------------------

func sendUserSignupMail(firstname, email string) error {
	sendMyselfNewUserMail(firstname, email)

	// TODO(remy): implement me!

	return nil
}

func sendMyselfNewUserMail(firstname, email string) error {
	if !UseMail {
		return nil
	}

	host := fmt.Sprintf("%s:%d", config.Config.SmtpHost, config.Config.SmtpPort)

	auth := auth()
	buff := bytes.Buffer{}

	// headers
	mailHeader(&buff, "me@remy.io", "New user!")

	// content
	buff.WriteString(fmt.Sprintf("<h1>New User</h1><ul><li>Email: %s</li><li>Firstname: %s</li></ul>\r\n", email, firstname))
	buff.WriteString("\r\n")

	// send
	err := smtp.SendMail(host, auth, Sender, []string{"me@remy.io"}, buff.Bytes())
	if err != nil {
		return log.Err("sendNewUserMail", err)
	}

	return nil
}
