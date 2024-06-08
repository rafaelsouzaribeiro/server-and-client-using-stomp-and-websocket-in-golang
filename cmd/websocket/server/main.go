package main

import "github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/websocket/server"

func main() {

	svc := server.NewServer("localhost", "ws", 8080)
	svc.ServerWebsocket()

}
