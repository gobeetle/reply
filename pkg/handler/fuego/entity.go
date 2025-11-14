package fuego

import (
	"github.com/gobeetle/reply/internal/decoder"
	"github.com/gobeetle/reply/internal/iface"
	resppkg "github.com/gobeetle/reply/internal/response"
	"github.com/gobeetle/reply/internal/transform"
)

type FuegoReplyHandler struct {
	transform.ResponseTransformOpt
	decoder iface.ResponseDecoder
}

func NewFuegoHandler() *FuegoReplyHandler {
	return &FuegoReplyHandler{
		decoder: decoder.NewDefaultDecoder(),
	}
}

func (r *FuegoReplyHandler) WithResponseTransformOpt(
	opt transform.ResponseTransformOpt,
) *FuegoReplyHandler {
	r.ResponseTransformOpt = opt
	return r
}

func (r *FuegoReplyHandler) WithResponseDecoder(
	decoder iface.ResponseDecoder,
) *FuegoReplyHandler {
	if decoder == nil {
		return r
	}
	r.decoder = decoder
	return r
}

func (r *FuegoReplyHandler) Decode(obj any) (iface.ErrorCoder, error) {
	if r.decoder == nil {
		return nil, nil
	}
	base, err := r.decoder.Decode(obj)
	if err != nil {
		return nil, err
	}
	return r.handle(base)
}

func (r *FuegoReplyHandler) handle(base resppkg.Response) (iface.ErrorCoder, error) {
	// first transform the response
	data, err := r.ResponseTransformOpt.Transform(base)
	if err != nil {
		return nil, err
	}

	return data, nil
}
