package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/websocket/client"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/infra/web/websocket/server"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
	"github.com/stretchr/testify/assert"
)

func TestSystemMessages(t *testing.T) {
	svc := server.NewServer("localhost", "/ws", 8080)
	go svc.ServerWebsocket()

	time.Sleep(1 * time.Second)

	con := client.NewClient("localhost", "ws", 8080)
	channel := make(chan dto.Payload)
	go con.ClientWebsocket("Client 1", "Hello 1", channel)

	var messages []dto.Payload

	timeout := time.After(5 * time.Second)
loop:
	for {
		select {
		case msg := <-channel:
			messages = append(messages, msg)
		case <-timeout:
			break loop
		}
	}

	for _, msg := range messages {
		if msg.Username == "Info" {
			assert.Contains(t, msg.Message, "User Client 1 connected")
		} else if msg.Username == "Client 1" {
			assert.Equal(t, "Client 1", msg.Username)
			assert.Contains(t, msg.Message, "Hello 1")
		}
	}

}

func BenchmarkWriter(b *testing.B) {
	svc := server.NewServer("localhost", "/ws", 8080)
	go svc.ServerWebsocket()

	time.Sleep(1 * time.Second)

	con := client.NewClient("localhost", "ws", 8080)

	for i := 0; i < b.N; i++ {
		channel := make(chan dto.Payload)
		go con.ClientWebsocket(fmt.Sprintf("Client %d", i), fmt.Sprintf("Hello %d", i), channel)

		select {
		case obj := <-channel:
			fmt.Printf("%s: %s\n", obj.Username, obj.Message)
		case <-time.After(5 * time.Second):
			b.Fatal("Timeout waiting for message")
		}

	}
}
