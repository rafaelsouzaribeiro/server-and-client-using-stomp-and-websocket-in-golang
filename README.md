Server and client using STOMP and WebSocket in Go. WebSocket with notifications for logged-in and logged-out users with log server, implementing <a href="https://github.com/rafaelsouzaribeiro/jwt-auth" title="JWT authentication">JWT authentication</a> and STOMP authentication.

For multiple messages on the websocket, the username is linked to a connection, so whenever you send more than one message to the username, you need to use client.Connect()

 ```go
package main

import (
	"fmt"
	"time"

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
		for range time.Tick(time.Second * 1) {
			client3.Send("Client 3", "Hello 3")
		}
	}()

	go func() {
		client4.Connect()
		go client4.Listen()
		for range time.Tick(time.Second * 1) {
			client4.Send("Client 4", "Hello 4")
		}
	}()

	for msg := range channel {
		fmt.Printf("%s: %s\n", msg.Username, msg.Message)
	}
}


```
send messages to a single username

 ```go
package main

import (
	"fmt"

	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/websocket/client"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
)

func main() {

	channel := make(chan dto.Payload)
	
	go func() {
		client := client.NewClient("localhost", "ws", 8080)
		defer client.Conn.Close()
		client.Connect()
		client.ClientWebsocket("Client 1", "Hello 1", channel)
	}()

	for obj := range channel {
		fmt.Printf("%s: %s\n", obj.Username, obj.Message)
	}

	close(channel)

}




 ```