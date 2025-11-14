package data

import (
	"encoding/json"

	"github.com/gobeetle/reply/internal/envelope"
	"github.com/gobeetle/reply/internal/iface"
)

type DataReply struct {
	code    int
	data    any
	message []string
}

func New(data any) *DataReply {
	return &DataReply{
		data: data,
	}
}

func (d *DataReply) WithCode(code int) *DataReply {
	d.code = code
	return d
}

func (d *DataReply) WithMessage(message ...string) *DataReply {
	d.message = message
	return d
}

func (d *DataReply) Data() any {
	return d.data
}

func (d *DataReply) StatusCode() int {
	return d.code
}

func (d *DataReply) Message() []string {
	return d.message
}

func (d *DataReply) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		envelope.ReplyEnvelope{
			Code:    d.code,
			Data:    d.data,
			Message: d.message,
		},
	)
}

var (
	_ iface.DataCoder     = (*DataReply)(nil)
	_ iface.MessageCoder  = (*DataReply)(nil)
	_ iface.JsonMarshaler = (*DataReply)(nil)
)
