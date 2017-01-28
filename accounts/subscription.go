package accounts

import (
	"time"

	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"
)

type Plan struct {
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
	Basic Plan = Plan{
		Name:        "Basic",
		Price:       500,
		Duration:    time.Hour * 24 * 90,
		DurationStr: "3 months",
	}
	Starter Plan = Plan{
		Name:        "Starter",
		Price:       900,
		Duration:    time.Hour * 24 * 180,
		DurationStr: "6 months",
	}
	Year Plan = Plan{
		Name:        "Year",
		Price:       1700,
		Duration:    time.Hour * 24 * 360,
		DurationStr: "1 year",
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
		("uid", "owner_uid", "stripe_customer_token", "stripe_charge_token", "plan", "price", "end", "stripe_response", "creation_time", "last_update")
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, uid, u.Uid, u.StripeToken, chargeId, plan.Name, plan.Price, end, string(json), t, t); err != nil {
		return err
	}

	return nil
}

func TrialInfos(uid uuid.UUID) (bool, time.Time, error) {
	// TODO(remy): implement me
	return true, time.Now().Add(time.Hour * 24), nil
}

func SubscriptionInfos(uid uuid.UUID) (bool, Plan, time.Time, error) {
	// TODO(remy): implement me
	return false, Basic, time.Now().Add(time.Hour * 24), nil
}
