package server

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	jwtauth "github.com/rafaelsouzaribeiro/jwt-auth/pkg/middleware"
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
var messageConnnected = make((map[string]bool))
var messageExists = make((map[*websocket.Conn]bool))
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

		if verify(msg.Username, &messageConnnected) {
			fmt.Printf("User connected: %s\n", msg.Username)

			systemMessag := dto.Payload{
				Username: fmt.Sprintf("Info %s", msg.Username),
				Message:  fmt.Sprintf("User %s connected", msg.Username),
			}

			err := msg.Conn.WriteJSON(systemMessag)

			if err != nil {
				fmt.Println("Error sending system message:", err)
				msg.Conn.Close()
				deleteUserByUserName(msg.Username, false)
			}
		}

		fmt.Printf("%s : %s \n", msg.Username, msg.Message)

		err := msg.Conn.WriteJSON(msg)

		if err != nil {
			fmt.Println("Error sending system message:", err)
			msg.Conn.Close()
			deleteUserByUserName(msg.Username, false)
		}

	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	defer func() {
		username := getUsernameByConnection(conn)

		mu.Lock()
		delete(messageExists, conn)
		mu.Unlock()

		if username != "" {
			_ = verifyToken(username, authHeader, conn)

			fmt.Printf("User %s disconnected\n", username)
			deleteUserByUserName(username, true)

			mu.Lock()
			delete(messageConnnected, username)
			mu.Unlock()

			conn.Close()
		}

	}()

	for {

		var msgs dto.Payload
		err := conn.ReadJSON(&msgs)
		if err != nil {
			break
		}

		if !verifyExistsUser(msgs.Username, conn) {
			if verifyCon(conn, &messageExists) {
				systemMessag := dto.Payload{
					Username: "info",
					Message:  fmt.Sprintf("User already exists: %s", msgs.Username),
				}

				fmt.Printf("User already exists: %s\n", msgs.Username)
				deleteUserByConn(conn, false)

				conn.WriteJSON(systemMessag)
			}
			continue
		}

		if !verifyToken(msgs.Username, authHeader, conn) {
			continue
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
		msgs.Conn = conn
		broadcast <- msgs

	}
}

func verifyToken(username, authHeader string, conn *websocket.Conn) bool {
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	cre, errs := jwtauth.NewCredential(3600, "rafael1234", nil)

	if errs != nil {
		fmt.Printf("Error jwt auth: %s", errs)
		return false
	}

	if cre.TokenExpired(tokenString) {
		systemMessag := dto.Payload{
			Username: "info",
			Message:  fmt.Sprintln("Sorry, your token expired"),
		}

		fmt.Printf("%s: Sorry, your token expired\n", username)
		deleteUserByConn(conn, false)
		conn.WriteJSON(systemMessag)
		conn.Close()
		return false
	}

	return true
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

func deleteUserByConn(conn *websocket.Conn, close bool) {
	mu.Lock()
	defer mu.Unlock()
	for k, user := range users {
		if user.conn == conn {
			if close {
				user.conn.Close()
			}
			delete(users, k)
		}
	}
}

func verifyExistsUser(u string, conn *websocket.Conn) bool {
	mu.Lock()
	defer mu.Unlock()
	for _, user := range users {
		if user.conn != conn && u == user.username {
			return false
		}

	}

	return true
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

func verifyCon(s *websocket.Conn, variable *map[*websocket.Conn]bool) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := (*variable)[s]; !exists {
		(*variable)[s] = true
		return true
	}
	return false

}
