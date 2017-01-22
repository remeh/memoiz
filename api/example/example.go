package example

import (
	"net/http"

	"remy.io/memoiz/api"
)

type Example struct {
}

func (c Example) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.RenderOk(w)
}
