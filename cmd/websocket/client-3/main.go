package main

import (
	"fmt"

	"github.com/rafaelsouzaribeiro/websocket-and-stomp-client-server-in-golang/internal/infra/web/websocket/client"
	"github.com/rafaelsouzaribeiro/websocket-and-stomp-client-server-in-golang/internal/usecase/dto"
)

func main() {

	con := client.NewClient("localhost", "ws", 8080)
	channel := make(chan dto.Payload)
	con.ClientWebsocket("Client 3", "Hello 3", channel)

	for obj := range channel {
		fmt.Printf("%s: %s\n", obj.Username, obj.Message)
	}
}
