package accounts

import (
	"net/http"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/api"
)

type PasswordReset struct{}

type pwdResetBody struct {
	Password string `json:"pwd"`
	Token    string `json:"token"`
}

func (c PasswordReset) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// parse parameters
	// ----------------------

	var body pwdResetBody
	api.ReadJsonBody(r, &body)

	if len(body.Password) == 0 || len(body.Token) == 0 {
		api.RenderBadParameters(w)
		return
	}

	// TODO(remy): validate password

	// updates password with the email with this token
	// ----------------------

	var pwdc string
	var err error

	if pwdc, err = accounts.Crypt(body.Password); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	if _, err := accounts.DAO().PwdReset(body.Token, pwdc); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderOk(w)
}
