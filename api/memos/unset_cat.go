package memos

import (
	"net/http"
	"time"

	"remy.io/memoiz/api"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/uuid"

	"github.com/gorilla/mux"
)

type UnsetCat struct {
}

func (c UnsetCat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	// parameters
	// ----------------------

	vars := mux.Vars(r)
	muid, err := uuid.Parse(vars["uid"])
	if err != nil {
		api.RenderBadParameters(w)
		return
	}

	if err := memos.DAO().UnsetCat(uid, muid, time.Now()); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderOk(w)
}
