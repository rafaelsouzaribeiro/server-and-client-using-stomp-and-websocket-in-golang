package main

import (
	"fmt"

	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/websocket/client"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
)

func main() {
	channel := make(chan dto.Payload)

	client3 := client.NewClient("localhost", "ws", 8080)
	client3.Channel = channel

	client4 := client.NewClient("localhost", "ws", 8080)
	client4.Channel = channel

	go func() {
		client3.Connect()
		go client3.Listen()
		client3.Send("Client 3", "Hello 3.1")
		client3.Send("Client 3", "Hello 3.2")
	}()

	go func() {
		client4.Connect()
		go client4.Listen()
		client4.Send("Client 4", "Hello 4.1")
		client4.Send("Client 4", "Hello 4.2")
	}()

	for msg := range channel {
		fmt.Printf("%s: %s\n", msg.Username, msg.Message)
	}
}
