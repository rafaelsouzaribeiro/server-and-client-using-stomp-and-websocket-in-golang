package payload

import (
	"github.com/go-stomp/stomp/v3/frame"
)

type Payload struct {
	Destination string
	Message     string
	Header      *frame.Header
	ContentType string
}
