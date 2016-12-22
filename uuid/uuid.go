// UUID
//
// Rémy Mathieu © 2016

package uuid

import (
	"fmt"

	"github.com/pborman/uuid"
)

type UUID uuid.UUID

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// ----------------------

func Parse(s string) (UUID, error) {
	r := uuid.Parse(s)
	if r == nil {
		return nil, fmt.Errorf("can't read the uuid: %s", s)
	}
	return UUID(r), nil
}

func Equal(left, right UUID) bool {
	return uuid.Equal(uuid.UUID(left), uuid.UUID(right))
}

// IsNil returns whether the given id should
// be considered nil or not.
func IsNil(id UUID) bool {
	if id == nil {
		return true
	}
	if len(id) == 0 {
		return true
	}
	if uuid.Equal(uuid.UUID(id), uuid.NIL) {
		return true
	}

	return false
}

// New returns a new random UUID.
func New() UUID {
	return UUID(uuid.Parse(uuid.New()))
}

func (u *UUID) Scan(value interface{}) error {
	s, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("UUID must be scanned from string")
	}

	// parse the value

	if v, err := Parse(string(s)); err != nil {
		return err
	} else {
		*u = v
		return nil
	}
}

// json
// ----------------------

func (u UUID) MarshalJSON() ([]byte, error) {
	return uuid.UUID(u).MarshalJSON()
}

// NOTE(remy): we'll probably need to implements UnmarshalJSON
