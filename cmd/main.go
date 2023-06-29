package main

import (
	"log"

	"github.com/fisher60/chessh/pkg/connection_handler"
	"github.com/gliderlabs/ssh"
)

func main() {
	ssh.Handle(connection_handler.ConnectionHandler)

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe("0.0.0.0:2222", nil, ssh.HostKeyFile("./.ssh/id_ed25519")))
}
