// Adapter to check that the user is
// correctly authed.
//
// Rémy Mathieu © 2016

package adapter

import (
	"net/http"

	"remy.io/scratche/api"
)

type AuthHandler struct {
	handler http.Handler
}

func (a AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)
	if uid.IsNil() {
		api.RenderForbiddenJson(w)
		return
	}

	a.handler.ServeHTTP(w, r)
}

// AuthAdapter creates a route which will force testing the
// auth cookie.
func AuthAdapter(handler http.Handler) http.Handler {
	return AuthHandler{
		handler: handler,
	}
}
