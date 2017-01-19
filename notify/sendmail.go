package notify

import (
	"bytes"
	"fmt"
	"net/smtp"

	"remy.io/scratche/config"
	"remy.io/scratche/log"
	"remy.io/scratche/notify/template"
)

var (
	UseMail = false

	Sender = "scratche@remy.io"
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

func Prout() {
	println("a!")
}

type Temp struct {
	Category string
	Count    int
}

func Sendmail() {
	host := fmt.Sprintf("%s:%d", config.Config.SmtpHost, config.Config.SmtpPort)

	body := []byte("To: me@remy.io\r\nSubject: Hello!\r\n\r\nCorpus\r\n")

	auth := smtp.PlainAuth("",
		config.Config.SmtpLogin,
		config.Config.SmtpPassword,
		config.Config.SmtpHost,
	)
	fmt.Println(host, body, auth)

	buff := bytes.Buffer{}

	t := Temp{
		Category: "Movie",
		Count:    1,
	}

	html := template.Root.Lookup("base.html")
	html.Execute(&buff, t)

	fmt.Println(string(buff.Bytes()))
	/*
		err := smtp.SendMail(host, auth, Sender, []string{"me@remy.io"}, body)
		if err != nil {
			log.Error(err)
		}
	*/
}
