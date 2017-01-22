package accounts

import (
	"net/http"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/api"
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
