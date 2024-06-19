package client

import (
	"fmt"
	"log"

	"github.com/go-stomp/stomp/v3"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
)

type Client struct {
	host     string
	port     int
	username string
	passcode string
}

func NewClient(host string, port int, username, passcode string) *Client {
	return &Client{
		host:     host,
		port:     port,
		username: username,
		passcode: passcode,
	}
}

func (c *Client) Send(pay *dto.Payload, channel chan<- dto.Payload) {
	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(c.username, c.passcode),
	}

	conn, err := stomp.Dial("tcp", fmt.Sprintf("%s:%d", c.host, c.port), options...)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer conn.Disconnect()

	sub, err := conn.Subscribe(pay.Destination, stomp.AckAuto)
	if err != nil {
		log.Fatalf("Error subscribing to destination: %v", err)
	}
	defer sub.Unsubscribe()

	err = conn.Send(pay.Destination, "text/plain", []byte(pay.Message), nil)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	for msg := range sub.C {
		channel <- dto.Payload{
			Message:     string(msg.Body),
			Destination: msg.Destination,
			Header:      msg.Header,
			ContentType: msg.ContentType,
		}
	}
}
