package dto

import (
	"github.com/go-stomp/stomp/v3/frame"
	"github.com/gorilla/websocket"
)

type Payload struct {
	Destination string
	Message     string `json:"message"`
	Header      *frame.Header
	ContentType string
	Username    string `json:"username"`
	Id          string
	Conn        *websocket.Conn
}
