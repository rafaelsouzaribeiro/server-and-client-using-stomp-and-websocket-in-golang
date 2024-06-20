package client

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	jwtauth "github.com/rafaelsouzaribeiro/jwt-auth/pkg/middleware"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
)

type Client struct {
	host    string
	port    int
	pattern string
	Conn    *websocket.Conn
	Channel chan dto.Payload
	Token   string
}

func NewClient(host, pattern string, port int) *Client {
	return &Client{
		host:    host,
		port:    port,
		pattern: pattern,
	}
}

func (client *Client) Connect() {
	header := client.GenerateToken()
	url := fmt.Sprintf("ws://%s:%d/%s", client.host, client.port, client.pattern)
	conn, _, err := websocket.DefaultDialer.Dial(url, *header)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	client.Conn = conn
}

func (client *Client) ClientWebsocket(username, message string, channel chan<- dto.Payload) {
	errs := client.Conn.WriteJSON(dto.Payload{Username: username, Message: message})
	if errs != nil {
		log.Println("Error writing message:", errs)
		return
	}

	for {
		var msg dto.Payload
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		channel <- msg
	}
}

func (client *Client) Send(username, message string) {
	errs := client.Conn.WriteJSON(dto.Payload{Username: username, Message: message})
	if errs != nil {
		log.Println("Error writing message:", errs)
		return
	}
}

func (client *Client) Listen() {
	defer close(client.Channel)
	for {
		var msg dto.Payload
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		client.Channel <- msg
	}
}

func (client *Client) GenerateToken() *http.Header {
	header := http.Header{}

	cre, err := jwtauth.NewCredential(3600, "rafael1234", nil)

	if err != nil {
		fmt.Printf("Error jwt auth: %s", err)
	}

	claims := map[string]interface{}{}
	token, errs := cre.CreateToken(claims)

	if errs != nil {

		fmt.Printf("Error create jwt auth: %s", errs)
	}
	header.Add("Authorization", "Bearer "+token)

	return &header
}
