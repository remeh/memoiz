package cards

import (
	"net/http"
	"time"

	"remy.io/scratche/api"
	"remy.io/scratche/cards"
	"remy.io/scratche/uuid"
)

type Put struct {
}

func (c Put) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	// TODO(remy): auth

	var body struct {
		CardUid uuid.UUID `json:"card_uid"`
		Text    string    `json:"text"`
	}

	api.ReadJsonBody(r, &body)
	sc, err := cards.DAO().UpdateText(uid, body.CardUid, body.Text, time.Now())
	if err != nil {
		api.RenderErrJson(w, err)
		return
	}
	api.RenderJson(w, 200, sc)
}
