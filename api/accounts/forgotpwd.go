package accounts

import (
	"net/http"
	"time"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/api"
	"remy.io/memoiz/log"
	"remy.io/memoiz/notify/email"
)

var (
	PwdResetTokenValidity time.Duration = time.Hour * 4
)

type ForgotPassword struct{}

type forgotPwdBody struct {
	Email string `json:"email"`
}

func (c ForgotPassword) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// parse parameters
	// ----------------------

	var body forgotPwdBody
	if err := api.ReadJsonBody(r, &body); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	if len(body.Email) == 0 {
		api.RenderBadParameters(w)
		return
	}

	// find user uid
	// ----------------------

	var err error
	var su accounts.SimpleUser

	if su, _, err = accounts.DAO().UserByEmail(body.Email); err != nil {
		log.Error("ForgotPassword:", err)
		return
	}

	if su.Uid == nil {
		api.RenderOk(w)
		return
	}

	// updates user token and validity time
	// ----------------------

	tok := accounts.PasswordResetToken(su.Uid)
	validUntil := time.Now().Add(PwdResetTokenValidity)

	var updated bool

	if updated, err = accounts.DAO().UpdatePwdResetToken(su.Uid, tok, validUntil); err != nil {
		log.Error("ForgotPassword:", err)
	}

	if updated {
		// sends the user the email to reset its password
		go func() {
			if err := email.SendPasswordResetMail(su, tok); err != nil {
				log.Error("ForgotPassword:", err)
			}
		}()
	}

	api.RenderOk(w)
}
