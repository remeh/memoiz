package accounts

import (
	"net/http"
	"time"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/api"
)

type Check struct{}

func (c Check) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the user
	// ----------------------

	t := api.ReadSessionToken(r)
	if len(t) == 0 {
		api.RenderForbiddenJson(w)
		return
	}

	var s accounts.Session
	var exists bool

	if s, exists = accounts.GetSession(t); !exists {
		api.RenderForbiddenJson(w)
		return
	} else {
		// TODO(remy): refresh the user session last hit?
	}

	// check the subscription of this user
	// in his session.
	// ----------------------

	if s.ValidUntil.Before(time.Now()) {
		// not valid anymore
		api.RenderPaymentRequired(w, s.Plan)
		return
	}

	api.RenderOk(w)
}
