package dto

import (
	"github.com/go-stomp/stomp/v3/frame"
)

type Payload struct {
	Destination string
	Message     string `json:"message"`
	Header      *frame.Header
	ContentType string
	Username    string `json:"username"`
	Id          string
}
