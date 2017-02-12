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

	// TODO(remy): validate password
	if len(body.Password) == 0 {
		api.RenderBadParameter(w, "password")
		return
	}

	if len(body.Token) == 0 {
		api.RenderBadParameter(w, "token")
		return
	}

	// updates password with the email with this token
	// ----------------------

	var pwdc string
	var err error

	if pwdc, err = accounts.Crypt(body.Password); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	if done, err := accounts.DAO().PwdReset(body.Token, pwdc); err != nil {
		api.RenderErrJson(w, err)
		return
	} else if !done {
		api.RenderBadParameter(w, "token")
		return
	}

	api.RenderOk(w)
}
