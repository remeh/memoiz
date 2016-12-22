// TVGame Backend
//
// Main
//
// Rémy Mathieu © 2016

package main

import (
	"net/http"

	"remy.io/scratche/api/adapter"
	"remy.io/scratche/api/example"
	"remy.io/scratche/config"
	l "remy.io/scratche/log"
)

func main() {
	l.Info("starting the runtime.")

	s := NewServer()

	declareApiRoutes(s)
	//startJobs()

	l.Info("listening on", config.Config.ListenAddr)

	err := s.Start()
	if err != nil {
		l.Error(err.Error())
	}
}

func log(h http.Handler) http.Handler {
	return adapter.LogAdapter(h)
}

func declareApiRoutes(s *Server) {
	s.AddApi("/example", log(example.Example{}))

	//// User routes
	//// ----------------------

	//rt.AddApi("/1.0/start", log(user.Start{}), "POST")

	//// Schedule
	//// ----------------------

	//rt.AddApi("/1.0/schedule", log(schedule.Schedule{}), "GET")

	//// Auction routes
	//// ----------------------

	//rt.AddApi("/1.0/bid/{auction_id}", log(auction.Bid{}), "POST")

	//rt.AddApi("/1.0/market/sell", log(auction.Sell{}), "POST")

	//rt.AddApi("/1.0/market/show/{auction_id}", log(auction.Info{Type: "show"}), "GET")
	//rt.AddApi("/1.0/market/host/{auction_id}", log(auction.Info{Type: "host"}), "GET")

	//// Channel routes
	//// ----------------------

	//rt.AddApi("/1.0/channel", log(channel.Create{}), "POST")
	//rt.AddApi("/1.0/channel/name", log(channel.RandomName{}), "GET")
}
