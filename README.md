Server and client using STOMP and WebSocket. WebSocket provides notifications for logged-in and logged-out users, and STOMP authentication.

For multiple messages the username is linked to a connection so whenever you send more than one message to the username you need to use client.Connect()

 ```go
package main

import (
	"fmt"

	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/websocket/client"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
)

func main() {
	client := client.NewClient("localhost", "ws", 8080)
	client.Connect()
	client.Channel = make(chan dto.Payload)
	defer client.Conn.Close()

	go client.Listen()

	go func() {
		client.Send("Client 3", "Hello 3")
		client.Send("Client 3", "Hello 4")
	}()

	for obj := range client.Channel {
		fmt.Printf("%s: %s\n", obj.Username, obj.Message)
	}

}
```