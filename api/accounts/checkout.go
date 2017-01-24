package accounts

import (
	"fmt"
	"net/http"

	"remy.io/memoiz/api"
)

type Checkout struct{}

type checkoutBody struct {
	Card      checkoutCard `json:"card"`
	ClientIp  string       `json:"client_ip"`
	CreatedTs int          `json:"created"`
	LiveMode  bool         `json:"livemode"`
	Tok       string       `json:"id"`
}

type checkoutCard struct {
	Id             string `json:"id"`
	Brand          string `json:"brand"`
	AddressZip     string `json:"address_zip"`
	AddressCity    string `json:"address_city"`
	AddressCountry string `json:"address_country"`
	AddressState   string `json:"address_state"`
	Country        string `json:"country"`
	ExpYear        int    `json:"exp_year"`
	ExpMonth       int    `json:"exp_month"`
	Funding        string `json:"funding"`
	Last4          string `json:"last4"`
}

func (c Checkout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO(remy): ensure the user is logged in

	var body checkoutBody
	if err := api.ReadJsonBody(r, &body); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	if len(body.Id) == 0 || strings.HasPrefix(body.Id, "cus_") {
		api.RenderBadParameters(w)
		return
	}

	// TODO(remy): store customer token in the user account
	// TODO(remy): store the checkout info in a custom table with the full JSON
	// TODO(remy): proceed to do the payment with Stripe
}
