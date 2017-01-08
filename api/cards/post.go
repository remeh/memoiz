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
		Enrich  bool      `json:"enrich"`
	}

	api.ReadJsonBody(r, &body)

	var err error
	var sc cards.Card

	if body.CardUid.IsNil() {
		sc, err = cards.DAO().New(uid, body.Text, time.Now())
	} else {
		sc, err = cards.DAO().UpdateText(uid, body.CardUid, body.Text, time.Now())
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
