// Memoiz Backend
//
// Main
//
// Rémy Mathieu © 2016

package main

import (
	"net/http"

	"remy.io/memoiz/api/accounts"
	"remy.io/memoiz/api/adapter"
	"remy.io/memoiz/api/example"
	"remy.io/memoiz/api/memos"
	"remy.io/memoiz/config"
	l "remy.io/memoiz/log"
)

func main() {
	l.Info("starting the runtime.")

	s := NewServer()

	declareApiRoutes(s)
	startJobs()

	l.Info("listening on", config.Config.ListenAddr)

	err := s.Start()
	if err != nil {
		l.Error(err.Error())
	}
}

func log(h http.Handler) http.Handler {
	return adapter.LogAdapter(h)
}

func auth(h http.Handler) http.Handler {
	return adapter.AuthAdapter(h)
}

func startJobs() {
}

func declareApiRoutes(s *Server) {
	s.AddApi("/example", log(example.Example{}))

	// Accounts
	// ----------------------

	s.AddApi("/1.0/plans", log(auth(accounts.Plans{})), "GET")

	s.AddApi("/1.0/accounts", log(accounts.Create{}), "POST")
	s.AddApi("/1.0/accounts", log(accounts.Check{}), "GET")
	s.AddApi("/1.0/accounts/login", log(accounts.Login{}), "POST")
	s.AddApi("/1.0/accounts/logout", log(accounts.Logout{}), "POST")
	s.AddApi("/1.0/accounts/checkout", log(auth(accounts.Checkout{})), "POST")

	s.AddApi("/1.0/emailing/unsubscribe/{token}", log(accounts.Unsubscribe{}), "GET")

	// Memos routes
	// ----------------------

	s.AddApi("/1.0/memos", log(auth(memos.Get{})), "GET")
	s.AddApi("/1.0/memos", log(auth(memos.Post{})), "POST")
	s.AddApi("/1.0/memos/switch/{left}/{right}", log(auth(memos.SwitchPosition{})), "POST")
	s.AddApi("/1.0/memos/{uid}/archive", log(auth(memos.Archive{})), "POST")
	s.AddApi("/1.0/memos/{uid}/rich", log(auth(memos.Rich{})), "GET")

}
