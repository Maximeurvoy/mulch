package mulch

import "time"

// TODO: add server timestamp

// SUCCESS & FAILURE will end a client connection (no?)
const (
	MessageSuccess = "SUCCESS"
	MessageFailure = "FAILURE"
)

const (
	MessageError   = "ERROR"
	MessageWarning = "WARNING"
	MessageInfo    = "INFO"
	MessageTrace   = "TRACE"
)

const MessageNoop = "NOOP"

const MessageNoTarget = ""

type Message struct {
	Time    time.Time `json:"time"`
	Type    string    `json:"type"`
	Target  string    `json:"target"`
	Message string    `json:"message"`
}

func NewMessage(mtype string, target string, message string) *Message {
	return &Message{
		Time:    time.Now(),
		Type:    mtype,
		Target:  target,
		Message: message,
	}
}
