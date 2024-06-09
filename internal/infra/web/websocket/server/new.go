package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rafaelsouzaribeiro/server-and-client-using-stomp-and-websocket-in-golang/internal/usecase/dto"
)

type Server struct {
	host    string
	port    int
	pattern string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan dto.Payload)
var messageBuffer []dto.Payload
var users = make(map[string]*websocket.Conn)

func NewServer(host, pattern string, port int) *Server {
	return &Server{
		host:    host,
		port:    port,
		pattern: pattern,
	}
}

func (server *Server) ServerWebsocket() {
	http.HandleFunc(server.pattern, handleConnections)

	go handleMessages()

	fmt.Printf("Server started on %s:%d \n", server.host, server.port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", server.host, server.port), nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}

func handleMessages() {
	for msg := range broadcast {
		messageBuffer = append(messageBuffer, msg)

		fmt.Printf("User connected: %s\n", msg.Username)

		for client := range clients {
			users[msg.Username] = client
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
				delete(users, msg.Username)
			}
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		username := getUsernameByConnection(conn)
		fmt.Printf("User %s disconnected\n", username)
		conn.Close()
		delete(clients, conn)
		delete(users, username)
	}()

	clients[conn] = true

	for _, msg := range messageBuffer {
		err := conn.WriteJSON(msg)
		if err != nil {
			delete(clients, conn)
			delete(users, msg.Username)
			fmt.Println(err)
			return
		}
	}

	for {
		var msg dto.Payload
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			return
		}

		broadcast <- msg
	}
}

func getUsernameByConnection(conn *websocket.Conn) string {
	for username, connection := range users {
		if connection == conn {
			return username
		}
	}
	return ""
}
