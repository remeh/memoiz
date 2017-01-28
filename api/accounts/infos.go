package accounts

import (
	"net/http"
	"time"

	"remy.io/memoiz/accounts"
	"remy.io/memoiz/api"
)

type Infos struct{}

type infosResp struct {
	Firstname              string    `json:"firstname"`
	Email                  string    `json:"email"`
	Trial                  bool      `json:"trial"`
	TrialValidUntil        time.Time `json:"free_trial_valid_until"`
	Subscribed             bool      `json:"subscribed"`
	Plan                   plan      `json:"plan"`
	SubscriptionValidUntil time.Time `json:"subscription_valid_until"`
}

func (c Infos) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := api.ReadUser(r)

	// gets the user information
	// ----------------------

	var su accounts.SimpleUser
	var err error

	if su, _, err = accounts.DAO().UserByUid(uid); err != nil {
		api.RenderErrJson(w, err)
		return
	}

	// gets the user sub if any
	// ----------------------

	resp := infosResp{
		Firstname: su.Firstname,
		Email:     su.Email,
	}

	trial, trialValidUntil, err := accounts.TrialInfos(uid)
	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	resp.Trial = trial
	resp.TrialValidUntil = trialValidUntil

	hasSub, plan, planValidUntil, err := accounts.SubscriptionInfos(uid)
	if err != nil {
		api.RenderErrJson(w, err)
		return
	}

	resp.Subscribed = hasSub
	resp.Plan = planToPlan(plan.Name, plan)
	resp.SubscriptionValidUntil = planValidUntil

	api.RenderJson(w, 200, resp)
}
