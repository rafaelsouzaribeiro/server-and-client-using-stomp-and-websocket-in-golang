package client

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/rafaelsouzaribeiro/websocket-and-stomp-client-server-in-golang/internal/usecase/dto"
)

type Client struct {
	host    string
	port    int
	pattern string
}

func NewClient(host, pattern string, port int) *Client {
	return &Client{
		host:    host,
		port:    port,
		pattern: pattern,
	}
}

func (client *Client) ClientWebsocket(username, message string, channel chan<- dto.Payload) {
	url := fmt.Sprintf("ws://%s:%d/%s", client.host, client.port, client.pattern)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	//defer conn.Close()

	errs := conn.WriteJSON(dto.Payload{Username: username, Message: message})
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

			channel <- msg
		}
	}()

}
