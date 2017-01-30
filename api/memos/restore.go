package memos

import (
	"net/http"
	"time"

	"remy.io/memoiz/api"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/uuid"

	"github.com/gorilla/mux"
)

type Restore struct{}

func (c Restore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	// read parameters
	// ----------------------

	var memoUid uuid.UUID
	var err error

	vars := mux.Vars(r)
	if memoUid, err = uuid.Parse(vars["uid"]); err != nil {
		api.RenderBadParameters(w)
		return
	}

	// test parameters
	// ----------------------

	if memoUid.IsNil() {
		api.RenderBadParameters(w)
		return
	}

	// ----------------------

	if err := memos.DAO().Restore(uid, memoUid, time.Now()); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderOk(w)
}
