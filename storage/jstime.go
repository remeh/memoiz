package storage

import (
	"database/sql/driver"
	"encoding/json"
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

func (js *JSTime) UnmarshalJSON(data []byte) error {
	var t int64
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	*js = JSTime(time.Unix(t/int64(1000), 0))
	return nil
}

func (js JSTime) Value() (driver.Value, error) {
	return driver.Value(time.Time(js)), nil
}
