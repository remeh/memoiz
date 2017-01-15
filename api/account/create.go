package account

import (
	"net/http"
	"time"

	"remy.io/scratche/account"
	"remy.io/scratche/api"
	"remy.io/scratche/uuid"
)

type Create struct{}

func (c Create) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	// TODO(remy): auth (do not permit user creation when logged in)

	// read parameters
	// ----------------------

	var body struct {
		Email     string
		Firstname string
		Password  string
	}

	if err := api.ReadJsonBody(r, &body); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	if len(body.Email) == 0 || account.ValidEmail(body.Email) {
		api.RenderBadParameter(w, "email")
		return
	}

	if len(body.Firstname) == 0 {
		api.RenderBadParameter(w, "firstname")
		return
	}

	if len(body.Password) == 0 ||
		account.IsPasswordSecure(body.Password) {
		api.RenderBadParameter(w, "password")
		return
	}

	// user creation
	// ----------------------

	uid := uuid.New()
	now := time.Now()

	if err := account.DAO().Create(uid, firstname, email, account.Crypt(password), now); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// resp
	// ----------------------

	resp := struct {
		api.Response
		Uid uuid.UUID
	}{
		Msg: "ok",
		Ok:  true,
		Uid: uid,
	}

	api.RenderJson(w, 200, resp)
}
