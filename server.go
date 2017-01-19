// Scratche backend.
//
// Listening server.
//
// Rémy Mathieu © 2016

package main

import (
	"net/http"

	"remy.io/scratche/config"
	l "remy.io/scratche/log"
	"remy.io/scratche/notify"
	"remy.io/scratche/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	s := &Server{}
	s.router = mux.NewRouter()
	return s
}

// Starts listening.
func (s *Server) Start() error {
	// Opens the database connection.
	l.Info("opening the database connection.")
	_, err := storage.Init(config.Config.ConnString)
	if err != nil {
		return err
	}

	// XXX(remy): remove
	notify.SendCategoryMail(notify.TmpGenerateCards())

	// Prepares the router serving the static pages and assets.
	s.prepareStaticRouter()

	// Handles static routes
	http.Handle("/", s.router)

	// Starts listening.
	err = http.ListenAndServe(config.Config.ListenAddr, nil)
	return err
}

// AddApi adds a route in the API router of the application.
func (s *Server) AddApi(pattern string, handler http.Handler, methods ...string) {
	s.router.PathPrefix("/api").Subrouter().Handle(pattern, handler).Methods(methods...)
}

// ----------------------

func (s *Server) prepareStaticRouter() {
	// Add the final route, the static assets and pages.
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir(config.Config.PublicDir)))
	l.Info("serving static from directory", config.Config.PublicDir)
}
