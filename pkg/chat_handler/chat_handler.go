package chat_handler

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/fisher60/chessh/pkg/menu_handler"
	"github.com/fisher60/chessh/pkg/session"
	"github.com/gliderlabs/ssh"
)

var messageQueue []string

func formatMessageQueue(mq []string) string {
	return strings.Join(mq, "\n")
}

func addMessageToQueue(m string) {
	messageQueue = append(messageQueue, m)

	for _, us := range menu_handler.SshSessions {
		renderTerminal(*us)
	}
}

func clearSession(s ssh.Session) {
	io.WriteString(s, "\033[H\033[2J")
}

func renderTerminal(us session.UserSession) {
	clearSession(us.Session)
	io.WriteString(us.Session, fmt.Sprintf("Chat -- %s\n\n", us.Session.User()))
	io.WriteString(us.Session, fmt.Sprintf("%s\n", formatMessageQueue(messageQueue)))
	io.WriteString(us.Session, fmt.Sprintf("> %s", us.Message))

}

func ChatHandler(s ssh.Session, us session.UserSession) {

	renderTerminal(us)

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
				// When enter key is pressed, send the message to the chat
				addMessageToQueue(fmt.Sprintf("%s: %s", s.User(), us.Message))
				us.Message = ""

			} else if key == 0x7f && len(us.Message) > 0 { // Backspace
				// When backspace key is pressed, delete the most recent character from message buffer
				us.Message = us.Message[:len(us.Message)-1]
			} else if key == 0x1b { // Escape
				// When escape key is pressed, exit the chat
				delete(menu_handler.SshSessions, us.Uuid)
				break MainLoop
			} else {
				// If the key did not match a special/defined character, add the character to message buffer
				us.Message = us.Message + string(key)
			}
			renderTerminal(us)
		}
	}

}
