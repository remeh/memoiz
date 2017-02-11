package accounts

import (
	"net/http"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/api"
)

type PasswordReset struct{}

func (c PasswordReset) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// parse parameters
	// ----------------------

	pwd := r.FormValue("pwd")
	token := r.FormValue("token")

	if len(pwd) == 0 || len(token) == 0 {
		api.RenderBadParameters(w)
		return
	}

	// TODO(remy): validate password

	// updates password with the email with this token
	// ----------------------

	var pwdc string
	var err error

	if pwdc, err = accounts.Crypt(pwd); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	println(pwdc)

	if err := accounts.DAO().PwdReset(token, pwdc); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	api.RenderOk(w)
}
