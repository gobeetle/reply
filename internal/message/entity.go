package message

import (
	"encoding/json"

	"github.com/gobeetle/reply/internal/envelope"
	"github.com/gobeetle/reply/internal/iface"
)

type MessageReply struct {
	code    int
	message []string
}

func New(message ...string) *MessageReply {
	return &MessageReply{
		message: message,
	}
}

func (m *MessageReply) WithCode(code int) *MessageReply {
	m.code = code
	return m
}

func (m *MessageReply) WithMessage(message ...string) *MessageReply {
	m.message = message
	return m
}

func (m *MessageReply) StatusCode() int {
	return m.code
}

func (m *MessageReply) Message() []string {
	return m.message
}

func (m *MessageReply) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		envelope.ReplyEnvelope{
			Code:    m.code,
			Message: m.message,
		},
	)
}

var (
	_ iface.MessageCoder  = (*MessageReply)(nil)
	_ iface.JsonMarshaler = (*MessageReply)(nil)
)
