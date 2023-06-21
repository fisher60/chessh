package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"
)

type userSession struct {
	Session ssh.Session
	Message string
}

func newUserSession(s ssh.Session, m string) *userSession {
	us := userSession{Session: s, Message: m}
	return &us
}

var messageQueue []string
var sshSessions = make(map[string]*userSession)

func formatMessageQueue(mq []string) string {
	return strings.Join(mq, "\n")
}

func addMessageToQueue(m string) {
	messageQueue = append(messageQueue, m)

	for _, us := range sshSessions {
		renderTerminal(*us)
	}
}

func clearSession(s ssh.Session) {
	io.WriteString(s, "\033[H\033[2J")
}

func renderTerminal(us userSession) {
	clearSession(us.Session)
	io.WriteString(us.Session, fmt.Sprintf("Chat -- %s\n\n", us.Session.User()))
	io.WriteString(us.Session, fmt.Sprintf("%s\n", formatMessageQueue(messageQueue)))
	io.WriteString(us.Session, fmt.Sprintf("> %s", us.Message))

}

func chessHandler(s ssh.Session) {
	userUUID := uuid.New().String()
	userSession := *newUserSession(s, "")
	sshSessions[string(userUUID)] = &userSession

	renderTerminal(userSession)

	reader := bufio.NewReader(s)
	keyPress := make(chan rune)

	go func() {
		for {
			r, _ := reader.ReadByte()
			keyPress <- rune(r)
		}
	}()

MainLoop:
	for {
		select {
		case key := <-keyPress:
			if key == 0xD { // Enter
				addMessageToQueue(fmt.Sprintf("%s: %s", s.User(), userSession.Message))
				userSession.Message = ""

			} else if key == 0x7f && len(userSession.Message) > 0 { // Backspace
				userSession.Message = userSession.Message[:len(userSession.Message)-1]
			} else if key == 0x1b {
				delete(sshSessions, string(userUUID))
				break MainLoop
			} else {
				userSession.Message = userSession.Message + string(key)
			}
			renderTerminal(userSession)
		}
	}

}

func main() {
	ssh.Handle(chessHandler)

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe("0.0.0.0:2222", nil, ssh.HostKeyFile("/.ssh/id_ed25519")))
}
