package accounts

import (
	"net/http"

	"remy.io/scratche/accounts"
	"remy.io/scratche/api"
)

type Logout struct{}

func (c Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := api.ReadSessionToken(r)
	if len(token) == 0 {
		api.RenderForbiddenJson(w)
		return
	}

	accounts.DeleteSession(token)
	api.RenderOk(w)
}
