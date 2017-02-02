package email

import (
	"bytes"
	"fmt"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/log"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/mind"
	"remy.io/memoiz/notify/template"
)

type semParam struct {
	SimpleUser accounts.SimpleUser
	Memos      enrichedMemos
}

type enrichedMemo struct {
	Memo        memos.Memo
	EnrichInfos mind.EnrichResult
}

type enrichedMemos []enrichedMemo

func SendEnrichedMemos(acc accounts.SimpleUser, ms memos.Memos, infos mind.EnrichResults) error {
	if !UseMail {
		return nil
	}

	if len(ms) == 0 {
		return fmt.Errorf("SendEnrichedMemos: called with 0 memos")
	}

	if len(ms) != len(infos) {
		return fmt.Errorf("SendEnrichedMemos: len(ms) != len(infos)")
	}

	//host := fmt.Sprintf("%s:%d", config.Config.SmtpHost, config.Config.SmtpPort)

	//auth := auth()
	buff := bytes.Buffer{}

	// headers
	mailHeader(&buff, acc.Email, "Hello!") // TODO(remy): generate a real title

	// content
	html := template.Root.Lookup("enriched_mail.html")
	if html == nil {
		return fmt.Errorf("SendEnrichedMail: can't find enriched template")
	}

	p := semParam{
		SimpleUser: acc,
		Memos:      buildEnrichedMemos(ms, infos),
	}

	if err := html.Execute(&buff, p); err != nil {
		log.Err("SendEnrichedMemos", err)
	}

	buff.WriteString("\r\n")

	dumpToFile(buff.Bytes())

	// send
	//err := smtp.SendMail(host, auth, Sender, []string{acc.Email}, buff.Bytes())
	//if err != nil {
	//	return log.Err("SendEnrichedMemos", err)
	//}

	return nil
}

func buildEnrichedMemos(ms memos.Memos, infos mind.EnrichResults) enrichedMemos {
	rv := make(enrichedMemos, len(ms))
	for i, m := range ms {
		rv[i].Memo = m
		rv[i].EnrichInfos = infos[i]
	}
	return rv
}
