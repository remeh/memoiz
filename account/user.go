package account

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"remy.io/scratche/uuid/"

	"golang.org/x/crypto/bcrypt"
)

func SetSessionCookie(w http.ResponseWriter, session db.Session) {
	cookie := &http.Cookie{
		Name:   "t",
		Value:  session.Token,
		MaxAge: 86400, // 1 day
	}
	http.SetCookie(w, cookie)
}

// IsPasswordSecure checks that the given password
// is strong enough to be used.
func IsPasswordSecure(password string) bool {
	// TODO(remy): check the password force
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
