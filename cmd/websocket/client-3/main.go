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
		con.ClientWebsocket("Client 3", "Hello 3", channel)
	}()

	for obj := range channel {
		fmt.Printf("%s: %s\n", obj.Username, obj.Message)
	}
}
