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
		renderChat(*us)
	}
}

func clearSession(s ssh.Session) {
	io.WriteString(s, "\033[H\033[2J")
}

func renderChat(us userSession) {
	clearSession(us.Session)
	io.WriteString(us.Session, fmt.Sprintf("Chat: As %s\n\n", us.Session.User()))
	io.WriteString(us.Session, fmt.Sprintf("%s\n", formatMessageQueue(messageQueue)))
	io.WriteString(us.Session, fmt.Sprintf("> %s", us.Message))

}

func buildMenuOptions(options map[int]string) string {
	var cache []string
	for k, v := range options {
		cache = append(cache, fmt.Sprintf("%d: %s\n", k, v))
	}

	return strings.Join(cache, "")
}

func renderMenu(us userSession, options map[int]string) {

	clearSession(us.Session)
	io.WriteString(us.Session, "Menu\n")
	io.WriteString(us.Session, buildMenuOptions(options))
}

func chatSession(us userSession) {

}

func sshHandler(s ssh.Session) {
	userUUID := uuid.New().String()
	us := *newUserSession(s, "")
	sshSessions[string(userUUID)] = &us

	var options = make(map[int]string)
	options[1] = "Chat"

	renderMenu(us, options)

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
				addMessageToQueue(fmt.Sprintf("%s: %s", s.User(), us.Message))
				us.Message = ""

			} else if key == 0x7f && len(us.Message) > 0 { // Backspace
				us.Message = us.Message[:len(us.Message)-1]
			} else if key == 0x1b { // escape
				delete(sshSessions, string(userUUID))
				break MainLoop
			} else {
				us.Message = us.Message + string(key)
			}
			renderChat(us)
		}
	}

}

func main() {
	ssh.Handle(sshHandler)

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe("0.0.0.0:2222", nil, ssh.HostKeyFile("./.ssh/id_ed25519")))
}
