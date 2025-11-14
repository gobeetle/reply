package reply

import "github.com/gobeetle/reply/internal/iface"

type (
	ErrorCoder   = iface.ErrorCoder
	StatusCoder  = iface.StatusCoder
	MessageCoder = iface.MessageCoder

	ErrorProvider   = iface.ErrorProvider
	MessageProvider = iface.MessageProvider

	Errors        = iface.Errors
	JsonMarshaler = iface.JsonMarshaler
)
