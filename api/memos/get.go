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

	// parse form
	// ----------------------

	r.ParseForm()
	s := r.Form.Get("s")

	memos.DAO()

	state := memos.MemoActive
	if s == "archived" {
		state = memos.MemoArchived
	}

	// get the memos
	// ----------------------

	cs, err := memos.DAO().GetByUser(uid, state)

	// render the response
	// ----------------------

	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderJson(w, 200, cs)
}
