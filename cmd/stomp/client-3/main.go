package main

import (
	"fmt"

	"github.com/rafaelsouzaribeiro/websocket-and-stomp-client-server-in-golang/pkg/payload"
	"github.com/rafaelsouzaribeiro/websocket-and-stomp-client-server-in-golang/pkg/stomp/client"
)

func main() {
	svc := client.NewClient("springboot", 8080, "admin", "1234")
	channel := make(chan payload.Payload)

	go svc.InitClient(&payload.Payload{
		Destination: "/topic/test",
		Message:     "Hello, STOMP 3!",
	}, channel)

	for cha := range channel {
		fmt.Printf("Message: %s Destination: %s Header: %s, ContentType: %s \n",
			cha.Message, cha.Destination, cha.Header, cha.ContentType)
	}
}
