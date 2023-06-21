package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"
)

type userSession struct {
	session  ssh.Session
	terminal *terminal.Terminal
}

func newUserSession(s ssh.Session, t *terminal.Terminal) *userSession {
	us := userSession{session: s, terminal: t}
	return &us
}

var messageQueue []string
var sshSessions = make(map[string]userSession)

func formatMessageQueue(mq []string) string {
	return strings.Join(mq, "\n")
}

func addMessageToQueue(m string) {
	messageQueue = append(messageQueue, m)

	for _, us := range sshSessions {
		renderTerminal(us)
	}
}

func clearSession(s ssh.Session) {
	io.WriteString(s, "\033[H\033[2J")
}

func renderTerminal(us userSession) {
	clearSession(us.session)
	io.WriteString(us.session, fmt.Sprintf("Chat -- %s\n\n", us.session.User()))
	io.WriteString(us.session, fmt.Sprintf("%s\n", formatMessageQueue(messageQueue)))

}

func chessHandler(s ssh.Session) {
	term := terminal.NewTerminal(s, "> ")

	userUUID := uuid.New().String()
	userSession := *newUserSession(s, term)
	sshSessions[string(userUUID)] = userSession

	renderTerminal(userSession)

	line := ""

	for {
		line, _ = userSession.terminal.ReadLine()
		if line == "quit" {
			break
		}

		addMessageToQueue(fmt.Sprintf("%s: %s", s.User(), line))
	}

	delete(sshSessions, string(userUUID))

}

func main() {
	ssh.Handle(chessHandler)

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe("localhost:2222", nil, ssh.HostKeyFile(".ssh/term_info_ed25519")))
}
