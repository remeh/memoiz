package cards

import (
	"net/http"
	"time"

	"remy.io/scratche/api"
	"remy.io/scratche/cards"
	"remy.io/scratche/mind"
	"remy.io/scratche/uuid"
)

type Post struct {
}

func (c Post) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	// TODO(remy): auth

	var body struct {
		CardUid uuid.UUID `json:"card_uid"`
		Text    string    `json:"text"`
	}

	api.ReadJsonBody(r, &body)

	var err error
	var sc cards.SimpleCard

	if uuid.IsNil(body.CardUid) {
		sc, err = cards.DAO().New(uid, body.Text, time.Now())
		go mind.Analyze(sc.Uid, body.Text)
	} else {
		sc, err = cards.DAO().UpdateText(body.CardUid, uid, body.Text, time.Now())
	}

	if err != nil {
		api.RenderErrJson(w, err)
		return
	}
	api.RenderJson(w, 200, sc)
}
