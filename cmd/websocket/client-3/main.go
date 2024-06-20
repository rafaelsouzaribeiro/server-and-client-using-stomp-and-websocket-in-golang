package main

import (
	"fmt"

	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/websocket/client"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
)

func main() {
	channel := make(chan dto.Payload)
	client := client.NewClient("localhost", "ws", 8080)
	client.Connect()
	client.Channel = channel
	defer client.Conn.Close()

	go client.Listen()

	go func() {
		client.Send("Client 3", "Hello 3")
		client.Send("Client 3", "Hello 4")
	}()

	for obj := range channel {
		fmt.Printf("%s: %s\n", obj.Username, obj.Message)
	}

}
