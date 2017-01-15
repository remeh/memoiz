package accounts

import (
	"net/http"
	"strings"
	"time"

	"remy.io/scratche/uuid"
)

const (
	SessionToken = "st"
)

type Session struct {
	// session token
	Token string
	// user uid
	Uid uuid.UUID

	CreationTime time.Time
	LastHit      time.Time
}

var sessions map[string]Session

func init() {
	sessions = make(map[string]Session)
}

// ----------------------

func SetSessionCookie(w http.ResponseWriter, s Session) {
	cookie := &http.Cookie{
		Name:   SessionToken,
		Value:  s.Token,
		Path:   "/api",
		MaxAge: 86400, // 1 day
	}
	http.SetCookie(w, cookie)
}

func NewSession(userUid uuid.UUID, t time.Time) Session {
	token := randSessTok()
	sessions[token] = Session{
		Token:        token,
		Uid:          userUid,
		CreationTime: t,
		LastHit:      t,
	}
	return sessions[token]
}

func RefreshSession(token string, t time.Time) {
	s, exists := sessions[token]
	if !exists {
		return
	}

	s.LastHit = t
	sessions[token] = s
}

func DeleteSession(token string) {
	delete(sessions, token)
}

func GetSession(token string) (Session, bool) {
	s, exists := sessions[token]
	return s, exists
}

// ----------------------

func randSessTok() string {
	rv := ""
	for i := 0; i < 3; i++ {
		str := uuid.New().String()
		str = strings.Replace(str, "-", "", -1)
		rv += str
	}
	return rv
}
