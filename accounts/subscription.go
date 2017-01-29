package accounts

import (
	"database/sql"
	"time"

	"remy.io/memoiz/log"
	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"
)

type Plan struct {
	Id string
	// Name of this plan.
	Name string
	// In cents, such as Stripe.
	Price uint64
	// Duration.
	Duration time.Duration
	// Duration description
	DurationStr string
}

var (
	// 7 days trial
	TrialDuration time.Duration = time.Hour * 24 * 7

	Basic Plan = Plan{
		Id:          "1",
		Name:        "Basic",
		Price:       500,
		Duration:    time.Hour * 24 * 90,
		DurationStr: "3 months",
	}
	Starter Plan = Plan{
		Id:          "2",
		Name:        "Starter",
		Price:       900,
		Duration:    time.Hour * 24 * 180,
		DurationStr: "6 months",
	}
	Year Plan = Plan{
		Id:          "3",
		Name:        "Year",
		Price:       1700,
		Duration:    time.Hour * 24 * 360,
		DurationStr: "1 year",
	}

	Plans map[string]Plan = map[string]Plan{
		"1": Basic,
		"2": Starter,
		"3": Year,
	}
)

// ----------------------

// AddSubscription stores in database the fact that an user
// has bought a Plan.
func AddSubscription(u SimpleUser, chargeId string, json []byte, t time.Time, plan Plan) error {
	uid := uuid.New()

	end := t.Add(plan.Duration)

	if _, err := storage.DB().Exec(`
		INSERT INTO "subscription"
		("uid", "owner_uid", "stripe_customer_token", "stripe_charge_token", "plan_id", "price", "end", "stripe_response", "creation_time", "last_update")
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, uid, u.Uid, u.StripeToken, chargeId, plan.Id, plan.Price, end, string(json), t, t); err != nil {
		return err
	}

	return nil
}

func TrialInfos(uid uuid.UUID) (bool, time.Time, error) {
	var t time.Time

	if err := storage.DB().QueryRow(`
		SELECT "creation_time" FROM "user"
		WHERE "uid" = $1
	`, uid).Scan(&t); err != nil {
		return false, t, log.Err("TrialInfos", err)
	}

	end := t.Add(TrialDuration)

	// still in trial
	if time.Now().Before(end) {
		return true, end, nil
	}

	// not anymore in trial
	return false, end, nil
}

func SubscriptionInfos(uid uuid.UUID) (bool, Plan, time.Time, error) {
	var end time.Time
	var planId string
	var plan Plan

	if err := storage.DB().QueryRow(`
		SELECT "end", "plan_id"
		FROM "subscription"
		WHERE
			owner_uid = $1
		ORDER BY "end" DESC
		LIMIT 1
	`, uid).Scan(&end, &planId); err != nil && err != sql.ErrNoRows {
		return false, Basic, end, log.Err("SubscriptionInfos", err)
	}

	var exists bool

	// has no plan
	if plan, exists = Plans[planId]; !exists {
		return false, Basic, time.Time{}, nil
	}

	// subscription still valid
	if time.Now().Before(end) {
		return true, plan, end, nil
	}

	// subscription not valid anymore
	return false, plan, end, nil
}
