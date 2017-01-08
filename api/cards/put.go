package cards

import (
	"net/http"
	"time"

	"remy.io/scratche/api"
	"remy.io/scratche/cards"
	"remy.io/scratche/mind"
	"remy.io/scratche/uuid"
)

type Put struct {
}

func (c Put) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO(remy): auth
	uid := api.ReadUser(r)

	// read parameters
	// ----------------------

	var body struct {
		CardUid uuid.UUID `json:"card_uid"`
		Text    string    `json:"text"`
		Enrich  bool      `json:"enrich"`
	}

	if err := api.ReadJsonBody(r, &body); err != nil {
		api.RenderBadParameters(w)
		return
	}

	// test parameter
	// ----------------------

	// TODO(remy): do we want to test text ?
	if body.CardUid.IsNil() {
		api.RenderBadParameters(w)
		return
	}

	sc, err := cards.DAO().UpdateText(uid, body.CardUid, body.Text, time.Now())
	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	if body.Enrich {
		go mind.Analyze(body.CardUid, body.Text)
	}

	api.RenderJson(w, 200, sc)
}
