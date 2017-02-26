package memos

import (
	"fmt"
	"net/http"
	"time"

	"remy.io/memoiz/api"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/mind"
	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"
)

type Post struct {
}

func (c Post) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	// ----------------------

	var body struct {
		MemoUid  uuid.UUID      `json:"memo_uid"`
		Text     string         `json:"text"`
		Enrich   bool           `json:"enrich"`
		Reminder storage.JSTime `json:"reminder"`
	}

	if err := api.ReadJsonBody(r, &body); err != nil {
		fmt.Println(err.Error())
		api.RenderBadParameters(w)
		return
	}

	var err error
	var sc memos.Memo

	if body.MemoUid.IsNil() {
		sc, err = memos.DAO().New(uid, body.Text, body.Reminder, time.Now())
	} else {
		sc, err = memos.DAO().UpdateText(uid, body.MemoUid, body.Text, body.Reminder, time.Now())
	}

	if body.Enrich {
		go mind.Analyze(sc.Uid, body.Text)
	}

	if err != nil {
		api.RenderErrJson(w, err)
		return
	}
	api.RenderJson(w, 200, sc)
}
