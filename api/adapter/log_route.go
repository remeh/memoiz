// Adapter to log a request on a route.
//
// Rémy Mathieu © 2016

package adapter

import (
	"fmt"
	"net/http"

	"remy.io/memoiz/config"
	"remy.io/memoiz/log"
)

type LogHandler struct {
	handler http.Handler
}

func (a LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", config.Config.AppUrl)
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// propagate to the next handler
	sWriter := &StatusWriter{w, 200}
	a.handler.ServeHTTP(sWriter, r)

	log.Info(fmt.Sprintf("Hit - %s %s %s referer[%s] user-agent[%s] addr[%s] code[%d]", r.Method, r.URL.String(), r.Proto, r.Referer(), r.UserAgent(), r.RemoteAddr, sWriter.Status))
}

// LogRoute creates a route which will log the route access.
func LogAdapter(handler http.Handler) http.Handler {
	return LogHandler{
		handler: handler,
	}
}

// ----------------------

type StatusWriter struct {
	http.ResponseWriter
	Status int
}

func (w *StatusWriter) WriteHeader(code int) {
	w.Status = code
	w.ResponseWriter.WriteHeader(code)
}
