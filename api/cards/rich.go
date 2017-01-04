package cards

import (
	"fmt"
	"net/http"

	"remy.io/scratche/api"
	"remy.io/scratche/cards"
	"remy.io/scratche/uuid"

	"github.com/gorilla/mux"
)

// Enrich returns the enriched information
// of the given card.
type Rich struct{}

func (c Rich) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO(remy): auth
	uid := api.ReadUser(r)

	// read parameters
	// ----------------------

	var cardUid uuid.UUID
	var err error

	vars := mux.Vars(r)

	if cardUid, err = uuid.Parse(vars["uid"]); err != nil {
		api.RenderBadParameters(w)
		return
	}

	// test parameter
	// ----------------------

	// TODO(remy): do we want to test text ?
	if cardUid.IsNil() {
		api.RenderBadParameters(w)
		return
	}

	// ----------------------

	ri, err := cards.DAO().GetRichInfo(uid, cardUid)

	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderJson(w, 200, ri)
}
