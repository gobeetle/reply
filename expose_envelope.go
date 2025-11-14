package reply

import (
	"github.com/gobeetle/reply/internal/envelope"
)

type (
	DataEnvelope[T any]        = envelope.DataEnvelope[T]
	ErrorEnvelope[E error]     = envelope.ErrorEnvelope[E]
	MessageEnvelope[M ~string] = envelope.MessageEnvelope[M]
	ReplyEnvelope              = envelope.ReplyEnvelope
)
