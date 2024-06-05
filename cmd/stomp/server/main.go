package main

import (
	svc "github.com/rafaelsouzaribeiro/server-websocket-and-stomp-golang/pkg/stomp/server"
)

func main() {

	server := svc.NewServer("springboot", 8080)
	server.InitServer()

}
