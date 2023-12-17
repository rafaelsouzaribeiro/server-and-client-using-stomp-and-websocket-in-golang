package main

import (
	"fmt"
	"log"

	"github.com/go-stomp/stomp"
)

func main() {
	// Endereço do servidor STOMP
	serverAddress := "localhost:61613"

	// Conectar ao servidor STOMP
	conn, err := stomp.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Disconnect()

	// Destino para enviar mensagens
	destination := "/topic/greetings"

	// Corpo da mensagem
	messageBody := "Hello, STOMP!"

	// Subscrever para receber mensagens
	sub, err := conn.Subscribe(destination, stomp.AckAuto)
	if err != nil {
		println(err)
	}
	defer sub.Unsubscribe()

	// Enviar uma mensagem
	err = conn.Send(destination, "text/plain", []byte(messageBody), nil)
	if err != nil {
		println(err)
	}

	// Aguardar por mensagens
	for {
		msg := <-sub.C
		fmt.Printf("Mensagem recebida: %s", msg.Body)
	}
}
