package cards

import (
	"net/http"

	"remy.io/scratche/api"
	"remy.io/scratche/cards"
)

type Get struct {
}

func (c Get) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cs, err := cards.DAO().GetByUser("", cards.CardActive)
	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderJson(w, 200, cs)
}
