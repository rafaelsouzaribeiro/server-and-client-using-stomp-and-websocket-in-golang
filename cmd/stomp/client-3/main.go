package main

import (
	"fmt"

	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/stomp/client"
	payload "github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
)

func main() {
	svc := client.NewClient("springboot", 8080, "admin", "1234")
	channel := make(chan payload.Payload)

	go svc.Send(&payload.Payload{
		Destination: "/topic/test",
		Message:     "Hello, STOMP 3!",
	}, channel)

	for cha := range channel {
		fmt.Printf("Message: %s Destination: %s Header: %s, ContentType: %s \n",
			cha.Message, cha.Destination, cha.Header, cha.ContentType)
	}
}
