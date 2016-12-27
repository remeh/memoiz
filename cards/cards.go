package cards

import (
	"database/sql/driver"

	"remy.io/scratche/uuid"
)

type CardState string

func (c CardState) String() string {
	return string(c)
}

func (c CardState) Value() (driver.Value, error) {
	return driver.Value(c.String()), nil
}

var (
	// CardActive is an active card of the user.
	CardActive CardState = "CardActive"
	// CardArchived has been archived by the user.
	CardArchived CardState = "CardArchived"
	// CardDeleted is used when the user has deleted the card.
	CardDeleted CardState = "CardDeleted"
)

// SimpleCard only contains necessary fields
// to represent a card.
type SimpleCard struct {
	Uid      uuid.UUID `json:"uid"`
	Text     string    `json:"text"`
	Position int       `json:"position"`
}
