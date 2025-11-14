package reply

import (
	"github.com/gobeetle/reply/internal/marshal"
)

type (
	ResponseMarshaler  = marshal.ResponseMarshaler
	ResponseMarshalOpt = marshal.ResponseMarshalOpt
)

// SetGlobalMarshaler sets a global marshaler
func SetGlobalMarshaler(m ResponseMarshaler) {
	marshal.SetGlobalMarshaler(m)
}

// GetGlobalMarshaler retrieves the global marshaler
func GetGlobalMarshaler() ResponseMarshaler {
	return marshal.GetGlobalMarshaler()
}
