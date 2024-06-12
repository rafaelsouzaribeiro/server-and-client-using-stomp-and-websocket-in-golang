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
var messageBufferMap = make(map[string][]dto.Payload)
var users = make(map[string]User)
var verifiedCon = make(map[string]bool)
var verifiedDes = make(map[string]bool)
var verifiedBuffer = make(map[string]bool)
var verifiedUser = make(map[string]bool)
var SystemMessage = make(map[string]bool)
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
		mu.Lock()
		messageBufferMap[msg.Id] = append(messageBufferMap[msg.Id], msg)
		mu.Unlock()

		if verify(msg.Username, &verifiedCon) {
			fmt.Printf("User connected: %s\n", msg.Username)
			mu.Lock()
			delete(verifiedDes, msg.Username)
			mu.Unlock()
			sendSystemMessage(fmt.Sprintf("User %s connected", msg.Username), msg.Username)
		}

		mu.Lock()
		for _, user := range users {
			if verifiedUser[msg.Id] {
				fmt.Printf("Duplicate ID found: %s\n", msg.Username)
				continue
			}
			verifiedUser[msg.Id] = true
			println(">>", user.username)
			err := user.conn.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
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
			sendSystemMessage(fmt.Sprintf("User %s disconnected", username), username)
		}

		deleteUserByUserName(username, true)
		conn.Close()
	}()

	mu.Lock()
	for _, msg := range messageBufferMap {
		for _, v := range msg {
			if verifiedBuffer[v.Id] {
				fmt.Printf("Duplicate ID found: %s\n", v.Id)
				continue
			}
			verifiedBuffer[v.Id] = true
			err := conn.WriteJSON(v)
			if err != nil {
				deleteUserByUserName(v.Username, false)
				fmt.Println(err)
				conn.Close()
				mu.Unlock()
				return
			}
		}
	}
	mu.Unlock()

	for {
		id := uuid.New().String()

		var msgs dto.Payload
		err := conn.ReadJSON(&msgs)
		if err != nil {
			break
		}

		mu.Lock()
		users[id] = User{
			conn:     conn,
			username: msgs.Username,
			id:       id,
		}
		verifiedBuffer = make(map[string]bool)
		verifiedUser = make(map[string]bool)
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

func sendSystemMessage(message, username string) {
	// systemMessage := dto.Payload{
	// 	Username: "Info",
	// 	Message:  message,
	// }

	// mu.Lock()
	// defer mu.Unlock()
	// for _, user := range users {
	// 	err := user.conn.WriteJSON(systemMessage)

	// 	if err != nil {
	// 		fmt.Println("Error sending system message:", err)
	// 		user.conn.Close()
	// 		deleteUserByUserName(user.username, false)
	// 	}
	// }
}
