package marshal

import (
	"encoding/json"

	"github.com/gobeetle/reply/internal/iface"
)

type ResponseMarshaler func(source iface.ErrorCoder) ([]byte, error)

type ResponseMarshalOpt struct {
	marshaler                ResponseMarshaler
	without_global_marshaler bool // default is false, meaning it will use the global marshaler if available
}

func (m *ResponseMarshalOpt) WithResponseMarshaler(marshaler ResponseMarshaler) *ResponseMarshalOpt {
	if marshaler == nil {
		return m
	}
	m.marshaler = marshaler
	return m
}

func (m *ResponseMarshalOpt) WithGlobalMarshaler(global_marshaler bool) *ResponseMarshalOpt {
	m.without_global_marshaler = !global_marshaler
	return m
}

func (m *ResponseMarshalOpt) Marshal(source iface.ErrorCoder) ([]byte, error) {
	if m.marshaler != nil {
		return m.marshaler(source)
	} else if !m.without_global_marshaler && GetGlobalMarshaler() != nil {
		return GetGlobalMarshaler()(source)
	}
	return json.Marshal(source)
}
