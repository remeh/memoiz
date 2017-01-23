package memos

import (
	"database/sql/driver"

	"remy.io/memoiz/mind"
	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"
)

type MemoState string

func (c MemoState) String() string {
	return string(c)
}

func (c MemoState) Value() (driver.Value, error) {
	return driver.Value(c.String()), nil
}

var (
	// MemoActive is an active memo of the user.
	MemoActive MemoState = "MemoActive"
	// MemoArchived has been archived by the user.
	MemoArchived MemoState = "MemoArchived"
	// MemoDeleted is used when the user has deleted the memo.
	MemoDeleted MemoState = "MemoDeleted"
)

type Memos []Memo

// Memo only contains necessary fields
// to represent a memo.
// RichInfo COULD be loaded.
type Memo struct {
	Uid      uuid.UUID `json:"uid"`
	Text     string    `json:"text"`
	Position int       `json:"-"`

	// NOTE(remy): everything in RichInfo is optional
	MemoRichInfo
}

type MemoRichInfo struct {
	// Loaded is true if the RichInfo are loaded.
	// Even partially.
	Loaded bool `json:"loaded,omitempty"`

	LastUpdate storage.JSTime `json:"last_update"` // timestamp
	Category   mind.Category  `json:"r_category,omitempty"`
	Image      string         `json:"r_img,omitempty"`
	Url        string         `json:"r_url,omitempty"`
	Title      string         `json:"r_title,omitempty"`
}

// ----------------------

// GroupByCategory regroups the slice of memos per Category.
func (cs Memos) GroupByCategory() map[mind.Category]Memos {
	rv := make(map[mind.Category]Memos)

	for _, c := range cs {
		v, exists := rv[c.MemoRichInfo.Category]

		if !exists {
			v = make(Memos, 0)
		}

		v = append(v, c)
		rv[c.MemoRichInfo.Category] = v
	}

	return rv
}
