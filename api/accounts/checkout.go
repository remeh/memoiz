package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/api"
	"remy.io/memoiz/log"
)

type Checkout struct{}

type checkoutBody struct {
	Card       checkoutCard `json:"card"`
	ClientIp   string       `json:"client_ip"`
	CreatedTs  int          `json:"created"`
	LiveMode   bool         `json:"livemode"`
	Tok        string       `json:"id"`
	Plan       string       `json:"plan"`
	CardHolder string       `json:"card_holder"`
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

type checkoutResp struct {
	api.Response
	Expiration time.Time `json:"expiration"`
}

var plans map[string]accounts.Plan = map[string]accounts.Plan{
	"1": accounts.Basic,
	"2": accounts.Starter,
	"3": accounts.Year,
}

func (c Checkout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

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
			Email: acc.Email,
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

	var plan accounts.Plan

	switch body.Plan {
	case "1", "2", "3":
	default:
		api.RenderBadParameter(w, "plan")
		return
	}

	plan = plans[body.Plan]

	// TODO(remy): if an error occurred here, we should retry 1 or 2 times.

	// proceed to the payment
	// ----------------------

	params := &stripe.ChargeParams{
		Amount:   plan.Price,
		Currency: "eur",
		Desc:     fmt.Sprintf("%s %s exp: %s", plan.Name, acc.Email, now.Add(plan.Duration)),
	}
	params.Meta = map[string]string{
		"CardHolder": body.CardHolder,
		"Email":      acc.Email,
	}

	params.Customer = acc.StripeToken

	charge, err := charge.New(params)
	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// re-render the json to store it in the subscription
	// ----------------------

	data, err := json.Marshal(body)
	if err != nil {
		log.Error("Checkout: while re-rendering the body:", err)
		data = []byte{}
	}

	// store the subscription
	// ----------------------

	if err := accounts.AddSubscription(acc, charge.ID, data, now, plan); err != nil {
		log.Error("Checkout: AddSubcription:", string(data))
		api.RenderErrJson(w, err)
		return
	}

	// response
	// ----------------------

	api.RenderOk(w)
}
