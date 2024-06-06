package main

import (
	svc "github.com/rafaelsouzaribeiro/websocket-and-stomp-client-server-in-golang/pkg/stomp/server"
)

func main() {

	server := svc.NewServer("springboot", 8080)
	server.InitServer()

}
