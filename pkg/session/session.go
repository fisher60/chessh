package session

import (
	"github.com/gliderlabs/ssh"
)

type UserSession struct {
	Session ssh.Session
	Message string
	Uuid    string
}

func NewUserSession(s ssh.Session, m string) *UserSession {
	us := UserSession{Session: s, Message: m}
	return &us
}
