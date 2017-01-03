package cards

import (
	"net/http"
	"time"

	"remy.io/scratche/api"
	"remy.io/scratche/cards"
	"remy.io/scratche/uuid"
)

type SwitchPosition struct {
}

func (c SwitchPosition) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO(remy): auth
	uid := api.ReadUser(r)

	// read parameters
	// ----------------------

	var body struct {
		LeftUid  uuid.UUID `json:"l"`
		RightUid uuid.UUID `json:"r"`
	}

	if err := api.ReadJsonBody(r, &body); err != nil {
		api.RenderBadParameters(w)
		return
	}

	// test parameters
	// ----------------------

	if body.LeftUid.IsNil() || body.RightUid.IsNil() {
		api.RenderBadParameters(w)
		return
	}

	// ----------------------

	if err := cards.DAO().SwitchPosition(body.LeftUid, body.RightUid, uid, time.Now()); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderOk(w)
}
