package main

import (
	svc "github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/stomp/server"
)

func main() {

	server := svc.NewServer("springboot", 8080)
	server.InitServer()

}
