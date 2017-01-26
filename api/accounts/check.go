package accounts

import (
	"net/http"

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

	if _, exists := accounts.GetSession(t); !exists {
		api.RenderForbiddenJson(w)
		return
	} else {
		// TODO(remy): refresh the user session last hit?
	}

	// ----------------------

	api.RenderOk(w)
}
