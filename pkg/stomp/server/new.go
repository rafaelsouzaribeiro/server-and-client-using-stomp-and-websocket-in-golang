package server

import (
	"fmt"
	"log"
	"net"

	"github.com/go-stomp/stomp/v3/server"
)

type Server struct {
	host string
	port int
}

func NewServer(host string, port int) *Server {

	return &Server{
		host: host,
		port: port,
	}
}

func (a Server) Authenticate(login, passcode string) bool {
	if login == "admin" && passcode == "1234" {
		return true
	}
	return false
}

func (b *Server) InitServer() {

	authenticator := Server{}

	stompServer := server.Server{
		Authenticator: authenticator,
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", b.host, b.port))

	if err != nil {
		log.Printf("Error starting server: %v", err)
	}

	log.Println("STOMP server running on port 8080...")

	if err := stompServer.Serve(ln); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
