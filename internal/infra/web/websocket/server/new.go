package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
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
	id       string
}

var broadcast = make(chan dto.Payload)
var users = make(map[string]User)
var verifiedCon = make(map[string]bool)
var verifiedDes = make(map[string]bool)
var mu sync.Mutex

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

		if verify(msg.Username, &verifiedCon) {
			fmt.Printf("User connected: %s\n", msg.Username)
			mu.Lock()
			delete(verifiedDes, msg.Username)
			mu.Unlock()

		}

		mu.Lock()

		for _, user := range users {
			if user.username != msg.Username {
				continue
			}

			systemMessag := dto.Payload{
				Username: fmt.Sprintf("Info %s", user.username),
				Message:  fmt.Sprintf("User %s connected", msg.Username),
			}

			fmt.Printf("%s : %s \n", user.username, msg.Message)
			err := user.conn.WriteJSON(systemMessag)

			if err != nil {
				fmt.Println("Error sending system message:", err)
				user.conn.Close()
				deleteUserByUserName(user.username, false)
			}

			err = user.conn.WriteJSON(msg)

			if err != nil {
				fmt.Println("Error sending system message:", err)
				user.conn.Close()
				deleteUserByUserName(user.username, false)
			}

		}
		mu.Unlock()
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
			mu.Lock()
			delete(verifiedCon, username)
			mu.Unlock()
		}

		if username != "" {
			deleteUserByUserName(username, true)
			conn.Close()
		}
	}()

	for {

		var msgs dto.Payload
		err := conn.ReadJSON(&msgs)
		if err != nil {
			break
		}

		mu.Lock()
		id := uuid.New().String()

		users[id] = User{
			conn:     conn,
			username: msgs.Username,
			id:       id,
		}
		mu.Unlock()

		msgs.Id = id
		broadcast <- msgs

	}
}

func getUsernameByConnection(conn *websocket.Conn) string {
	mu.Lock()
	defer mu.Unlock()
	for _, user := range users {
		if user.conn == conn {
			return user.username
		}
	}
	return ""
}

func deleteUserByUserName(username string, close bool) {
	mu.Lock()
	defer mu.Unlock()
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
	mu.Lock()
	defer mu.Unlock()
	if _, exists := (*variable)[s]; !exists {
		(*variable)[s] = true
		return true
	}
	return false
}
