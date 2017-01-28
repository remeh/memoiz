package accounts

import (
	"net/http"
	"strconv"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/api"
)

type Plans struct{}

type plansResp struct {
	Plans []plan `json:"plans"`
}

type plan struct {
	Name     string `json:"name"`
	Price    string `json:"price"`
	Duration string `json:"duration"`
}

func (c Plans) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps := []plan{
		planToPlan(accounts.Basic),
		planToPlan(accounts.Starter),
		planToPlan(accounts.Year),
	}

	resp := plansResp{
		Plans: ps,
	}

	api.RenderJson(w, 200, resp)

}

func planToPlan(p accounts.Plan) plan {
	round := strconv.Itoa(int(p.Price / 100))
	cents := strconv.Itoa(int(p.Price % 100))

	price := "$" + round
	if (p.Price % 100) > 0 {
		price += "," + cents
	}

	return plan{
		Name:     p.Name,
		Price:    price,
		Duration: p.DurationStr,
	}
}
