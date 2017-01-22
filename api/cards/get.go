package cards

import (
	"net/http"

	"remy.io/memoiz/api"
	"remy.io/memoiz/cards"
)

type Get struct {
}

func (c Get) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	cs, err := cards.DAO().GetByUser(uid, cards.CardActive)

	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderJson(w, 200, cs)
}
