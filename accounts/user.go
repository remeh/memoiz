package accounts

import (
	"remy.io/scratche/uuid"

	"golang.org/x/crypto/bcrypt"
)

type SimpleUser struct {
	Uid       uuid.UUID `json:"uid"`
	Firstname string    `json:"firstname"`
	Email     string    `json:"email"`
}

// IsPasswordSecure checks that the given password
// is strong enough to be used.
func IsPasswordSecure(password string) bool {
	// TODO(remy): check the password force
	return true
}

// ValidEmail returns whether the given email is
// valid or not.
func ValidEmail(email string) bool {
	// TODO(remy): valid email
	return true
}

// Crypt crypts the given password using bcrypt.
func Crypt(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(b), err
}

// Check validates that the hash is indeed derived from
// the given password.
func Check(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
