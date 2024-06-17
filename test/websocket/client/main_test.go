package main

import (
	"fmt"
	"sync"
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
	con.Connect()
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

var once sync.Once

func startServer() {
	svc := server.NewServer("localhost", "/ws", 8080)
	go svc.ServerWebsocket()
	time.Sleep(1 * time.Second)
}

var count int = 0

func BenchmarkWriter(b *testing.B) {
	once.Do(startServer)

	channel := make(chan dto.Payload)
	var messages []dto.Payload

	for i := 0; i < b.N; i++ {
		count++
		go func(count int) {
			client := client.NewClient("localhost", "ws", 8080)
			defer client.Conn.Close()
			client.Connect()
			client.ClientWebsocket(fmt.Sprintf("Client %d", count), fmt.Sprintf("Hello %d", count), channel)
		}(count)
	}

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

	for k, msg := range messages {
		c := fmt.Sprintf("Info %s %d: User %s %d connected", "Client ", k, "Client ", k)
		d := fmt.Sprintf("Client %d: Hello %d", k, k)

		if msg.Username == c {
			assert.Contains(b, msg.Message, d)
		}

	}
}
