package cards

import (
	"database/sql/driver"

	"remy.io/scratche/mind"
	"remy.io/scratche/storage"
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

type Cards []Card

// Card only contains necessary fields
// to represent a card.
// RichInfo COULD be loaded.
type Card struct {
	Uid      uuid.UUID `json:"uid"`
	Text     string    `json:"text"`
	Position int       `json:"-"`

	// NOTE(remy): everything in RichInfo is optional
	CardRichInfo
}

type CardRichInfo struct {
	// Loaded is true if the RichInfo are loaded.
	// Even partially.
	Loaded bool `json:"loaded,omitempty"`

	LastUpdate storage.JSTime `json:"last_update"` // timestamp
	Category   mind.Category  `json:"r_category,omitempty"`
	Image      string         `json:"r_img,omitempty"`
	Url        string         `json:"r_url,omitempty"`
	Title      string         `json:"r_title,omitempty"`
}
