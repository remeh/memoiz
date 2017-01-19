// Scratche Backend
//
// Main
//
// Rémy Mathieu © 2016

package main

import (
	"net/http"

	"remy.io/scratche/api/accounts"
	"remy.io/scratche/api/adapter"
	"remy.io/scratche/api/cards"
	"remy.io/scratche/api/example"
	"remy.io/scratche/config"
	l "remy.io/scratche/log"
	"remy.io/scratche/notify"
)

func main() {
	l.Info("starting the runtime.")

	s := NewServer()

	declareApiRoutes(s)
	startJobs()

	notify.Sendmail()

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

	s.AddApi("/1.0/accounts", log(accounts.Create{}), "POST")
	s.AddApi("/1.0/accounts/login", log(accounts.Login{}), "POST")
	s.AddApi("/1.0/accounts/logout", log(accounts.Logout{}), "POST")

	// Cards routes
	// ----------------------

	s.AddApi("/1.0/cards", log(auth(cards.Get{})), "GET")
	s.AddApi("/1.0/cards", log(auth(cards.Post{})), "POST")
	s.AddApi("/1.0/cards/switch/{left}/{right}", log(auth(cards.SwitchPosition{})), "POST")
	s.AddApi("/1.0/cards/{uid}/archive", log(auth(cards.Archive{})), "POST")
	s.AddApi("/1.0/cards/{uid}/rich", log(auth(cards.Rich{})), "GET")

}
