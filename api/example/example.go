package example

import (
	"net/http"

	"remy.io/scratche/api"
)

type Example struct {
}

func (c Example) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.RenderOk(w)
}
