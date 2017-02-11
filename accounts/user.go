package accounts

import (
	"strings"

	"remy.io/memoiz/uuid"

	"golang.org/x/crypto/bcrypt"
)

type SimpleUser struct {
	Uid       uuid.UUID `json:"uid"`
	Firstname string    `json:"firstname"`
	Email     string    `json:"email"`

	UnsubToken  string `json:"-"`
	StripeToken string `json:"-"`
}

// IsPasswordSecure checks that the given password
// is strong enough to be used.
func IsPasswordSecure(password string) bool {
	// TODO(remy): check the password force
	return true
}

// ValidEmail basically valids the given email.
func ValidEmail(email string) bool {
	return strings.Contains(email, "@") &&
		strings.Contains(email, ".")
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

// UnsubToken generates a random unsubscription
// token.
// It is composed of:
// - the char '1' (version)
// - first 8 chars of the user uid
// - with 3 randoms uuids (without -) appended.
func UnsubToken(uid uuid.UUID) string {
	end := randTok()
	start := uid.String()[0:8]
	return "1" + start + end
}

// Token generates a random password reset
// token.
// It is composed of:
// - the char '1' (version)
// - first 8 chars of the user uid
// - with 3 randoms uuids (without -) appended.
func PasswordResetToken(uid uuid.UUID) string {
	end := randTok()
	start := uid.String()[0:8]
	return "1" + start + end
}
