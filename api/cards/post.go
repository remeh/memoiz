package cards

import (
	"net/http"

	"remy.io/scratche/api"
	"remy.io/scratche/cards"
)

type Post struct {
}

func (c Post) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	// TODO(remy): auth

	var body struct {
		Text string `json:"text"`
	}

	api.ReadJsonBody(r, &body)
	sc, err := cards.DAO().New(uid, body.Text)
	if err != nil {
		api.RenderErrJson(w, err)
		return
	}
	api.RenderJson(w, 200, sc)
}
