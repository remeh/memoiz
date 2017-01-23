package memos

import (
	"net/http"

	"remy.io/memoiz/api"
	"remy.io/memoiz/memos"
	"remy.io/memoiz/uuid"

	"github.com/gorilla/mux"
)

// Enrich returns the enriched information
// of the given memo.
type Rich struct{}

func (c Rich) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// test parameter
	// ----------------------

	// TODO(remy): do we want to test text ?
	if memoUid.IsNil() {
		api.RenderBadParameters(w)
		return
	}

	// ----------------------

	ri, err := memos.DAO().GetRichInfo(uid, memoUid)

	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderJson(w, 200, ri)
}
