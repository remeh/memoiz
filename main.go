// Scratche Backend
//
// Main
//
// Rémy Mathieu © 2016

package main

import (
	"net/http"

	"remy.io/scratche/api/adapter"
	"remy.io/scratche/api/cards"
	"remy.io/scratche/api/example"
	"remy.io/scratche/config"
	l "remy.io/scratche/log"
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

func startJobs() {
}

func declareApiRoutes(s *Server) {
	s.AddApi("/example", log(example.Example{}))

	// Cards routes
	// ----------------------

	s.AddApi("/1.0/cards", log(cards.Get{}), "GET")
	s.AddApi("/1.0/cards", log(cards.Post{}), "POST")
	s.AddApi("/1.0/cards/switch", log(cards.SwitchPosition{}), "POST")
}
