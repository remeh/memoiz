package memos

import (
	"net/http"

	"remy.io/memoiz/api"
	"remy.io/memoiz/memos"
)

type Get struct {
}

func (c Get) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	cs, err := memos.DAO().GetByUser(uid, memos.MemoActive)

	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderJson(w, 200, cs)
}
