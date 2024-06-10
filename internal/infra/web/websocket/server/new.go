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

type User struct {
	conn     *websocket.Conn
	username string
	pointer  int
}

var broadcast = make(chan dto.Payload)
var messageBuffer []dto.Payload
var users = make(map[int]User)
var pointer = -1
var verifiedCon = make(map[string]bool)
var verifiedDes = make(map[string]bool)
var verifiedBuffer = make(map[string]bool)

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

		if verify(msg.Username, &verifiedCon) {
			fmt.Printf("User connected: %s\n", msg.Username)
			delete(verifiedDes, msg.Username)
			sendSystemMessage(fmt.Sprintf("User %s connected", msg.Username))

		}

		for _, user := range users {

			err := user.conn.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				user.conn.Close()
				deleteUserByUserName(user.username, false)
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

		if verify(username, &verifiedDes) {
			fmt.Printf("User %s disconnected\n", username)
			delete(verifiedCon, username)
			sendSystemMessage(fmt.Sprintf("User %s disconnected", username))
		}

		deleteUserByUserName(username, true)
		conn.Close()
	}()

	for _, msg := range messageBuffer {
		err := conn.WriteJSON(msg)
		if err != nil {
			deleteUserByUserName(msg.Username, false)
			fmt.Println(err)
			conn.Close()
			return
		}
	}

	pointer++

	for {
		var msgs dto.Payload
		err := conn.ReadJSON(&msgs)
		if err != nil {
			//fmt.Printf("Error reading message: %v\n", err)
			break
		}

		users[pointer] = User{
			conn:     conn,
			username: msgs.Username,
			pointer:  pointer,
		}

		broadcast <- msgs
	}
}

func getUsernameByConnection(conn *websocket.Conn) string {
	for _, user := range users {
		if user.conn == conn {
			return user.username
		}
	}
	return ""
}

func deleteUserByUserName(username string, close bool) {
	for k, user := range users {
		if user.username == username {
			if close {
				user.conn.Close()
			}

			delete(users, k)
		}
	}
}

func verify(s string, variable *map[string]bool) bool {
	if !(*variable)[s] {
		(*variable)[s] = true
		return true
	}
	return false
}

func sendSystemMessage(message string) {
	systemMessage := dto.Payload{
		Username: "Info",
		Message:  message,
	}

	for _, user := range users {
		err := user.conn.WriteJSON(systemMessage)

		if verify(message, &verifiedBuffer) {
			messageBuffer = append(messageBuffer, systemMessage)
		}

		if err != nil {
			fmt.Println("Error sending system message:", err)
			user.conn.Close()
			deleteUserByUserName(user.username, false)
		}
	}
}
