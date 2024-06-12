package main

import (
	"fmt"

	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/websocket/client"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
)

func main() {

	channel := make(chan dto.Payload)

	for i := 0; i < 2; i++ {
		go func(i int) {
			client := client.NewClient("localhost", "ws", 8080)
			client.Connect()
			client.ClientWebsocket(fmt.Sprintf("Client %d", i), fmt.Sprintf("Hello %d", i), channel)
			client.Conn.Close()
		}(i)
	}

	for obj := range channel {
		fmt.Printf("%s: %s\n", obj.Username, obj.Message)
	}

	close(channel)

}
