package connection_handler

import (
	"github.com/fisher60/chessh/pkg/chat_handler"
	"github.com/fisher60/chessh/pkg/menu_handler"
	"github.com/fisher60/chessh/pkg/session"
	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"
)

func ConnectionHandler(s ssh.Session) {
	userUUID := uuid.New().String()
	userSession := *session.NewUserSession(s, "")
	menu_handler.SshSessions[string(userUUID)] = &userSession

	chat_handler.ChatHandler(userSession)
}
