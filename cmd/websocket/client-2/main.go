package main

import (
	"fmt"

	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/websocket/client"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
)

func main() {

	channel2 := make(chan dto.Payload)

	for i := 2; i < 4; i++ {
		go func(i int) {
			client := client.NewClient("localhost", "ws", 8080)
			defer client.Conn.Close()
			client.Connect()
			client.ClientWebsocket(fmt.Sprintf("Client %d", i), fmt.Sprintf("Hello %d", i), channel2)
		}(i)
	}

	for objs := range channel2 {
		fmt.Printf("%s: %s\n", objs.Username, objs.Message)
	}

	close(channel2)
}
