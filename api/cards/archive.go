package cards

import (
	"net/http"
	"time"

	"remy.io/memoiz/api"
	"remy.io/memoiz/cards"
	"remy.io/memoiz/uuid"

	"github.com/gorilla/mux"
)

type Archive struct {
}

func (c Archive) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// test parameters
	// ----------------------

	if cardUid.IsNil() {
		api.RenderBadParameters(w)
		return
	}

	// ----------------------

	if err := cards.DAO().Archive(uid, cardUid, time.Now()); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderOk(w)
}
