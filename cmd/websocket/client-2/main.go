package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/rafaelsouzaribeiro/websocket-and-stomp-client-server-in-golang/internal/usecase/dto"
)

func main() {
	url := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	errs := conn.WriteJSON(dto.Payload{Username: "client 2", Message: "Hello 2"})
	if errs != nil {
		log.Println("Error writing message:", errs)
		return
	}

	go func() {
		for {
			var msg dto.Payload
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			fmt.Printf("%s: %s\n", msg.Username, msg.Message)
		}
	}()

	select {}

}
