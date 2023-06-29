package main

import (
	"log"

	"github.com/fisher60/chessh/pkg/chat_handler"
	"github.com/gliderlabs/ssh"
)

func main() {
	ssh.Handle(chat_handler.ChatHandler)

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe("0.0.0.0:2222", nil, ssh.HostKeyFile("/.ssh/id_ed25519")))
}
