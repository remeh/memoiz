package cards

import "remy.io/scratche/uuid"

// SimpleCard only contains necessary fields
// to represent a card.
type SimpleCard struct {
	Uid      uuid.UUID `json:"uid"`
	Text     string    `json:"text"`
	Position int       `json:"position"`
}
