package main

import "github.com/rafaelsouzaribeiro/websocket-and-stomp-client-server-in-golang/internal/infra/web/websocket/server"

func main() {

	svc := server.NewServer("localhost", "ws", 8080)
	svc.ServerWebsocket()

}
