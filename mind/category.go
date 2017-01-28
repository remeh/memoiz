package mind

import (
	"database/sql/driver"
	"fmt"
)

//go:generate stringer -type=Category

type Category int64

type Categories []Category

const (
	Uncategorized Category = iota
	Artist                 // 1
	Actor                  // 2
	Book                   // 3
	Date                   // 4
	Movie                  // 5
	Music                  // 6
	Person                 // 7
	Place                  // 8
	Serie                  // 9
	Video                  // 10
	VideoGame              // 11
	Food                   // 12
)

func (c *Category) Scan(src interface{}) error {
	var i int64
	var ok bool

	if i, ok = src.(int64); !ok {
		return fmt.Errorf("Category must be read from int")
	}

	*c = Category(i)
	return nil
}

func (c Category) Value() (driver.Value, error) {
	return driver.Value(int64(c)), nil
}

// json
// ----------------------

func (c Category) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, c.String())), nil
}

// TODO(remy): UnmarshalJSON
