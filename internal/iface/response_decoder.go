package iface

import "github.com/gobeetle/reply/internal/response"

type ResponseDecoder interface {
	Empty() (response.Response, error)
	Decode(obj any) (response.Response, error)
}
