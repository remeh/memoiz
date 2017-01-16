package cards

import (
	"net/http"
	"time"

	"remy.io/scratche/api"
	"remy.io/scratche/cards"
	"remy.io/scratche/uuid"

	"github.com/gorilla/mux"
)

type SwitchPosition struct {
}

func (c SwitchPosition) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	// read parameters
	// ----------------------

	var err error
	var leftUid, rightUid uuid.UUID

	vars := mux.Vars(r)

	if leftUid, err = uuid.Parse(vars["left"]); err != nil {
		api.RenderBadParameters(w)
		return
	}
	if rightUid, err = uuid.Parse(vars["right"]); err != nil {
		api.RenderBadParameters(w)
		return
	}

	if leftUid.IsNil() || rightUid.IsNil() {
		api.RenderBadParameters(w)
		return
	}

	// ----------------------

	if err := cards.DAO().SwitchPosition(leftUid, rightUid, uid, time.Now()); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderOk(w)
}
