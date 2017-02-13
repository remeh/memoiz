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

type semParam struct {
	SimpleUser accounts.SimpleUser
	Memos      enrichedMemos
}

type enrichedMemo struct {
	Memo        memos.Memo
	EnrichInfos mind.EnrichResult
}

type enrichedMemos []enrichedMemo

// ----------------------

type ImageFirst enrichedMemos

func (ms ImageFirst) Len() int          { return len(ms) }
func (ms ImageFirst) Swap(i int, j int) { ms[i], ms[j] = ms[j], ms[i] }
func (ms ImageFirst) Less(i int, j int) bool {
	ii := ms[i].EnrichInfos
	ij := ms[j].EnrichInfos

	if len(ii.ImageUrl) != 0 && len(ij.ImageUrl) == 0 {
		return true
	}

	return false
}

// SendEnrichedMemos sends to the given user the list of memos
// enriched by the given infos.
// It also stores the email in the dumpDir directory using
// the sendUid as filename.
func SendEnrichedMemos(acc accounts.SimpleUser, ms memos.Memos, infos mind.EnrichResults, dumpDir string, sendUid uuid.UUID) error {
	if !UseMail {
		return nil
	}

	if len(ms) == 0 {
		return fmt.Errorf("SendEnrichedMemos: called with 0 memos")
	}

	if len(ms) != len(infos) {
		return fmt.Errorf("SendEnrichedMemos: len(ms) != len(infos)")
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
	mailHeader(&buff, acc.Email, buildTitle(em))

	// content
	html := template.Root.Lookup("enriched_mail.html")
	if html == nil {
		return fmt.Errorf("SendEnrichedMemos: can't find enriched template")
	}

	p := semParam{
		SimpleUser: acc,
		Memos:      em,
	}

	if err := html.Execute(&buff, p); err != nil {
		return log.Err("SendEnrichedMemos", err)
	}

	buff.WriteString("\r\n")

	dumpToFile(dumpDir, sendUid.String(), buff.Bytes())

	// send
	err := smtp.SendMail(host(), auth(), Sender, []string{acc.Email}, buff.Bytes())
	if err != nil {
		return log.Err("SendEnrichedMemos", err)
	}

	return nil
}

func buildTitle(em enrichedMemos) string {
	str := ""
	for _, m := range em {
		if len(str) >= 70 {
			// avoid too long title
			break
		}

		if len(str) != 0 {
			str += " — "
		}

		if len(m.Memo.Title) != 0 {
			str += cutText(m.Memo.Title, 70)
		} else {
			str += cutText(m.Memo.Text, 70)
		}
	}

	if len(str) == 0 {
		return "Memos are waiting for you!"
	}

	return str
}

func cutText(str string, size int) string {
	if len(str) > size {
		return str[0:size] + "…"
	}
	return str
}

func buildEnrichedMemos(ms memos.Memos, infos mind.EnrichResults) enrichedMemos {
	rv := make(enrichedMemos, len(ms))
	for i, m := range ms {
		rv[i].Memo = m
		rv[i].EnrichInfos = infos[i]
	}
	return rv
}
