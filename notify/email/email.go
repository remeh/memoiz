package email

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"time"

	"remy.io/memoiz/config"
	"remy.io/memoiz/log"
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

func dumpToFile(b []byte) {
	ioutil.WriteFile(fmt.Sprintf("/tmp/tmp%d.html", time.Now().Unix()), b, 0644)
}

func mailHeader(buff *bytes.Buffer, email, title string) {
	buff.WriteString("MIME-version: 1.0;\r\n")
	buff.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n")
	buff.WriteString(fmt.Sprintf("To: %s\r\n", email))
	buff.WriteString(fmt.Sprintf("Subject: %s\r\n\r\n", title))
}

func auth() smtp.Auth {
	return smtp.PlainAuth("",
		config.Config.SmtpLogin,
		config.Config.SmtpPassword,
		config.Config.SmtpHost,
	)
}
