package storage

import (
	"fmt"
	"time"
)

// JSTimestamp is a basic timestamp
// but up to the ms and not to second.
// Because Javascript uses this unit...
type JSTime time.Time

func (js JSTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(js).Unix()*1000)), nil
}

// TODO(remy): marshal json
