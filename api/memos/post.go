package memos

import (
	"net/http"
	"time"

	"remy.io/memoiz/api"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/mind"
	"remy.io/memoiz/uuid"
)

type Post struct {
}

func (c Post) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	// ----------------------

	var body struct {
		MemoUid uuid.UUID `json:"memo_uid"`
		Text    string    `json:"text"`
		Enrich  bool      `json:"enrich"`
	}

	api.ReadJsonBody(r, &body)

	var err error
	var sc memos.Memo

	if body.MemoUid.IsNil() {
		sc, err = memos.DAO().New(uid, body.Text, time.Now())
	} else {
		sc, err = memos.DAO().UpdateText(uid, body.MemoUid, body.Text, time.Now())
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
