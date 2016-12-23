package cards

import "remy.io/scratche/uuid"

type CardState string

func (c CardState) String() string {
	return string(c)
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
