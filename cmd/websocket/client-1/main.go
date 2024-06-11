package main

import (
	"fmt"

	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/websocket/client"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
)

func main() {

	con := client.NewClient("localhost", "ws", 8080)
	con.Connect()
	channel := make(chan dto.Payload)
	go con.ClientWebsocket("Client 1", "Hello 1", channel)

	for obj := range channel {
		fmt.Printf("%s: %s\n", obj.Username, obj.Message)
	}

}
