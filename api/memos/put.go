package memos

import (
	"net/http"
	"time"

	"remy.io/memoiz/api"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/mind"
	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"
)

type Put struct {
}

func (c Put) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	// read parameters
	// ----------------------

	var body struct {
		MemoUid  uuid.UUID      `json:"memo_uid"`
		Text     string         `json:"text"`
		Enrich   bool           `json:"enrich"`
		Reminder storage.JSTime `json:"reminder"`
	}

	if err := api.ReadJsonBody(r, &body); err != nil {
		api.RenderBadParameters(w)
		return
	}

	// test parameter
	// ----------------------

	// TODO(remy): do we want to test text ?
	if body.MemoUid.IsNil() {
		api.RenderBadParameters(w)
		return
	}

	sc, err := memos.DAO().UpdateText(uid, body.MemoUid, body.Text, body.Reminder, time.Now())
	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	if body.Enrich {
		go mind.Analyze(body.MemoUid, body.Text)
	}

	api.RenderJson(w, 200, sc)
}
