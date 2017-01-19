package notify

import (
	"bytes"
	"fmt"
	"net/smtp"

	"remy.io/scratche/cards"
	"remy.io/scratche/config"
	"remy.io/scratche/log"
	"remy.io/scratche/mind"
	"remy.io/scratche/notify/template"
	"remy.io/scratche/uuid"
)

var (
	UseMail = false

	Sender = "ratch@remy.io"
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

func SendCategoryMail(cs map[mind.Category]cards.Cards) {
	host := fmt.Sprintf("%s:%d", config.Config.SmtpHost, config.Config.SmtpPort)

	auth := smtp.PlainAuth("",
		config.Config.SmtpLogin,
		config.Config.SmtpPassword,
		config.Config.SmtpHost,
	)

	buff := bytes.Buffer{}

	buff.WriteString("MIME-version: 1.0;\r\n")
	buff.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n")
	buff.WriteString("To: me@remy.io\r\n")
	buff.WriteString("Subject: Hello!\r\n\r\n")

	fmt.Println(cs)

	html := template.Root.Lookup("base.html")
	html.Execute(&buff, cs)

	buff.WriteString("\r\n")

	err := smtp.SendMail(host, auth, Sender, []string{"me@remy.io"}, buff.Bytes())
	if err != nil {
		log.Error(err)
	}
}

// XXX
func TmpGenerateCards() map[mind.Category]cards.Cards {
	uid, _ := uuid.Parse("12341234-1234-1234-1234-123412341234")
	cs, _ := cards.DAO().GetByUser(uid, cards.CardActive)

	rv := make(map[mind.Category]cards.Cards)

	for _, c := range cs {
		rv[mind.Movie] = append(rv[mind.Movie], c)
	}

	return rv
}
