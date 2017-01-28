package accounts

import (
	"net/http"
	"strconv"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/api"
)

type Plans struct{}

type plansResp struct {
	Plans map[string]plan `json:"plans"`
	Order []string        `json:"order"`
}

type plan struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Duration string `json:"duration"`
}

func (c Plans) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps := map[string]plan{
		"1": planToPlan("1", accounts.Basic),
		"2": planToPlan("2", accounts.Starter),
		"3": planToPlan("3", accounts.Year),
	}

	resp := plansResp{
		Plans: ps,
		Order: []string{"1", "2", "3"},
	}

	api.RenderJson(w, 200, resp)

}

func planToPlan(id string, p accounts.Plan) plan {
	round := strconv.Itoa(int(p.Price / 100))
	cents := strconv.Itoa(int(p.Price % 100))

	price := "$" + round
	if (p.Price % 100) > 0 {
		price += "," + cents
	}

	return plan{
		Id:       id,
		Name:     p.Name,
		Price:    price,
		Duration: p.DurationStr,
	}
}
