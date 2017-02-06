package accounts

import (
	"net/http"
	"strings"
	"time"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/api"
)

type Login struct{}

func (c Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// read parameters
	// ----------------------

	var body struct {
		Email    string
		Password string
	}

	if err := api.ReadJsonBody(r, &body); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	if len(body.Email) == 0 || !accounts.ValidEmail(body.Email) {
		api.RenderBadParameter(w, "email")
		return
	}

	if len(body.Password) == 0 {
		api.RenderBadParameter(w, "password")
		return
	}

	// gets user
	// ----------------------

	body.Email = strings.ToLower(body.Email)

	var su accounts.SimpleUser
	var hash string
	var err error
	now := time.Now()

	if su, hash, err = accounts.DAO().UserByEmail(body.Email); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	if !accounts.Check(hash, body.Password) {
		api.RenderForbiddenJson(w)
		return
	}

	// get the user subscription info
	// ----------------------

	var plan accounts.Plan
	var validUntil time.Time

	if _, plan, validUntil, err = accounts.SubscriptionInfos(su.Uid); err != nil {
		api.RenderForbiddenJson(w)
		return
	}

	// no plan, look for trial infos
	if plan == accounts.NoPlan {
		if _, validUntil, err = accounts.TrialInfos(su.Uid); err != nil {
			api.RenderForbiddenJson(w)
			return
		}
	}

	// create session
	// ----------------------

	s := accounts.NewSession(su.Uid, now, validUntil, plan)
	accounts.SetSessionCookie(w, s)
	api.RenderOk(w)
}
