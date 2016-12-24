package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"remy.io/scratche/log"
	"remy.io/scratche/uuid"
)

// Request
// ----------------------

// ReadUser reads the 'u' parameter in the request
// for a valid UUID.
func ReadUser(r *http.Request) uuid.UUID {
	r.ParseForm()

	param := r.Form.Get("u")
	if len(param) == 0 && len(param) != 36 {
		return nil
	}

	if id, err := uuid.Parse(param); err != nil {
		return nil
	} else {
		return id
	}
}

func ReadJsonBody(r *http.Request, object interface{}) error {
	if r == nil {
		return fmt.Errorf("ReadJsonBody: r == nil")
	}

	if object == nil {
		return fmt.Errorf("ReadJsonBody: object == nil")
	}

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return fmt.Errorf("ReadJsonBody: %s", err.Error())
	}

	return json.Unmarshal(data, object)
}

// Response
// ----------------------

// Response is the most basic type of response
// returned by this API.
// It only contains a message for the client.
type Response struct {
	Msg string `json:"msg"`
	Ok  bool   `json:"ok"`
}

// RenderJson renders in the given HTTP Response
// the marshaled object.
// If an error occured during the encoding, an error
// is logged and a 500 is returned.
func RenderJson(w http.ResponseWriter, code int, object interface{}) {
	data, err := json.Marshal(object)

	if err != nil {
		err := fmt.Errorf("while encoding in JSON an object of type '%s': %s",
			reflect.TypeOf(object).Name(), err.Error())
		RenderErrJson(w, err)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

// RenderErrJson renders a 500 http response
// with 'internal error' as response message.
// The error is also logged.
func RenderErrJson(w http.ResponseWriter, err error) {
	// log
	log.Error(err)

	// http response
	rdata, _ := json.Marshal(Response{
		Msg: "internal error",
		Ok:  false,
	})

	w.WriteHeader(500)
	w.Write(rdata)
	return
}

func RenderForbiddenJson(w http.ResponseWriter) {
	RenderBaseJson(w, 403, "forbidden")
}

func RenderBadParameters(w http.ResponseWriter) {
	RenderBaseJson(w, 400, "bad parameters")
}

func RenderOk(w http.ResponseWriter) {
	RenderBaseJson(w, 200, "ok")
}

// RenderBaseJson renders a simple JSON containing only
// a msg in the response.
func RenderBaseJson(w http.ResponseWriter, code int, msg string) {
	rdata, _ := json.Marshal(Response{
		Msg: msg,
		Ok:  code == 200,
	})

	w.WriteHeader(code)
	w.Write(rdata)
	return
}
