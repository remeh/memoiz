package accounts

import (
	"fmt"
	"net/http"
	"strings"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"

	"remy.io/memoiz/accounts"
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

// TODO(remy): deal with an invalid card

func (c Checkout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the user
	// ----------------------

	var acc accounts.SimpleUser

	uid := api.ReadUser(r)
	if uid == nil {
		api.RenderForbiddenJson(w)
		return
	}

	var err error

	if acc, _, err = accounts.DAO().UserByUid(uid); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	if acc.Uid.IsNil() {
		api.RenderForbiddenJson(w)
		return
	}

	// read the body
	// ----------------------

	var body checkoutBody
	if err := api.ReadJsonBody(r, &body); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	if len(body.Tok) == 0 || !strings.HasPrefix(body.Tok, "tok_") {
		api.RenderBadParameters(w)
		return
	}

	// if necessary, create a Stripe Token for this user.
	// ----------------------

	if len(acc.StripeToken) == 0 {
		// create a customer if not already created on Stripe
		// ----------------------

		customerParams := &stripe.CustomerParams{
			Desc: acc.Email,
		}
		customerParams.SetSource(body.Tok) // obtained with Stripe.js
		cus, err := customer.New(customerParams)

		if err != nil {
			api.RenderErrJson(w, err)
			return
		}

		// store the token for this user
		// ----------------------

		acc.StripeToken = cus.ID

		if err := accounts.DAO().UpdateStripeToken(acc); err != nil {
			api.RenderErrJson(w, err)
			return
		}
	}

	// proceed o do the payment
	// TODO(remy): handle that he have selected a plan.
	// ----------------------

	params := &stripe.ChargeParams{
		Amount:   500,
		Currency: "eur",
		Desc:     "Monthly example charge",
	}

	params.Customer = acc.StripeToken

	charge, err := charge.New(params)
	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// TODO(remy): store the checkout info in a custom table with the full JSON
	// TODO(remy): store the subscription state for this user

	fmt.Println(charge.Paid)
}
