package accounts

import (
	"net/http"
	"time"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/api"
	"remy.io/memoiz/log"
	"remy.io/memoiz/uuid"
)

var (
	PwdResetTokenValidity time.Duration = time.Hour * 4
)

type ForgotPassword struct{}

func (c ForgotPassword) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// parse parameters
	// ----------------------

	pwd := r.FormValue("email")

	// find user uid
	// ----------------------

	var err error
	var uid uuid.UUID

	if uid, err = accounts.DAO().UidByEmail(pwd); err != nil {
		log.Error("PasswordToken:", err)
		return
	}

	if uid == nil {
		api.RenderOk(w)
		return
	}

	// updates user token and validity time
	// ----------------------

	tok := accounts.PasswordResetToken(uid)
	validUntil := time.Now().Add(PwdResetTokenValidity)

	if err := accounts.DAO().UpdatePwdResetToken(uid, tok, validUntil); err != nil {
		log.Error("PasswordToken:", err)
	}

	api.RenderOk(w)
}
