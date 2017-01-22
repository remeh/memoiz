package accounts

import (
	"net/http"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/api"
	"remy.io/memoiz/log"

	"github.com/gorilla/mux"
)

// An user can click on a link in the email
// to unsubscribe himself from the mailing process.
type Unsubscribe struct{}

func (c Unsubscribe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// read parameters
	// ----------------------

	vars := mux.Vars(r)

	if len(vars["token"]) == 0 {
		api.RenderBadParameters(w)
		return
	}

	// insert this unsubscribe in the database
	// ----------------------

	if err := accounts.DAO().Unsubscribe(vars["token"], "email"); err != nil {
		log.Error("Unsubscribe:", err)
	}

	// atm, always return a 200
	// ----------------------
	w.Write([]byte("You're correctly unsubscribed."))
}
