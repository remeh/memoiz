package accounts

import (
	"net/http"
	"time"

	"remy.io/scratche/accounts"
	"remy.io/scratche/api"
	"remy.io/scratche/uuid"
)

type Create struct{}

func (c Create) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	if !uid.IsNil() {
		api.RenderBaseJson(w, 403, "you can't create many accounts")
		return
	}

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

	if len(body.Email) == 0 || !accounts.ValidEmail(body.Email) {
		api.RenderBadParameter(w, "email")
		return
	}

	if len(body.Firstname) == 0 {
		api.RenderBadParameter(w, "firstname")
		return
	}

	if len(body.Password) == 0 ||
		!accounts.IsPasswordSecure(body.Password) {
		api.RenderBadParameter(w, "password")
		return
	}

	// check for existence
	// ----------------------

	var err error
	if uid, err := accounts.DAO().UidByEmail(body.Email); err != nil {
		api.RenderErrJson(w, err)
		return
	} else {
		if !uid.IsNil() {
			api.RenderBaseJson(w, 409, "existing user")
			return
		}
	}

	// user creation
	// ----------------------

	var hash string

	uid = uuid.New()
	now := time.Now()

	if hash, err = accounts.Crypt(body.Password); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	if err := accounts.DAO().Create(uid, body.Firstname, body.Email, hash, now); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// resp
	// ----------------------

	resp := struct {
		api.Response
		Uid uuid.UUID `json:"uid"`
	}{
		Response: api.Response{
			Msg: "ok",
			Ok:  true,
		},
		Uid: uid,
	}

	api.RenderJson(w, 200, resp)
}
