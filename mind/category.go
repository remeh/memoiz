package mind

import (
	"database/sql/driver"
	"fmt"
)

type Category string

type Categories []Category

const (
	Uncategorized Category = "Uncategorized"
	Artist                 = "Artist"    // 1
	Actor                  = "Actor"     // 2
	Book                   = "Book"      // 3
	News                   = "News"      // 4
	Movie                  = "Movie"     // 5
	Music                  = "Music"     // 6
	Person                 = "Person"    // 7
	Place                  = "Place"     // 8
	Serie                  = "Series"    // 9
	Video                  = "Video"     // 10
	VideoGame              = "VideoGame" // 11
	Food                   = "Food"      // 12
)

func (c *Category) Scan(src interface{}) error {
	var i string
	var ok bool

	if i, ok = src.(string); !ok {
		return fmt.Errorf("Category must be read from int")
	}

	*c = Category(i)
	return nil
}

func (c Category) Value() (driver.Value, error) {
	return driver.Value(string(c)), nil
}

// json
// ----------------------

func (c Category) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, c)), nil
}

// TODO(remy): UnmarshalJSON
