package main

import (
	"fmt"

	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/websocket/client"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
)

func main() {
	channel := make(chan dto.Payload)

	go func() {
		con := client.NewClient("localhost", "ws", 8080)
		con.Connect()
		con.ClientWebsocket("Client 2", "Hello 2", channel)
		con.Conn.Close()
	}()

	for obj := range channel {
		fmt.Printf("%s: %s\n", obj.Username, obj.Message)
	}
}
