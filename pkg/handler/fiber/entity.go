package fiber

import (
	"github.com/gobeetle/reply/internal/decoder"
	"github.com/gobeetle/reply/internal/iface"
	"github.com/gobeetle/reply/internal/marshal"
	resppkg "github.com/gobeetle/reply/internal/response"
	"github.com/gobeetle/reply/internal/transform"
	"github.com/gofiber/fiber/v2"
)

type FiberReplyHandler struct {
	transform.ResponseTransformOpt
	marshal.ResponseMarshalOpt
	decoder iface.ResponseDecoder
	ctx     *fiber.Ctx
}

func NewFiberHandler(c *fiber.Ctx) *FiberReplyHandler {
	return &FiberReplyHandler{
		ctx:     c,
		decoder: decoder.NewDefaultDecoder(),
	}
}

func (r *FiberReplyHandler) WithResponseMarshalOpt(
	opt marshal.ResponseMarshalOpt,
) *FiberReplyHandler {
	r.ResponseMarshalOpt = opt
	return r
}

func (r *FiberReplyHandler) WithResponseTransformOpt(
	opt transform.ResponseTransformOpt,
) *FiberReplyHandler {
	r.ResponseTransformOpt = opt
	return r
}

func (r *FiberReplyHandler) WithResponseDecoder(
	decoder iface.ResponseDecoder,
) *FiberReplyHandler {
	if decoder == nil {
		return r
	}
	r.decoder = decoder
	return r
}

// this function returns an empty response, with status code set to StatusNoContent
func (r *FiberReplyHandler) Empty() error {
	return r.JSON(nil)
}

func (r *FiberReplyHandler) JSON(obj any) error {
	if r.decoder == nil {
		return nil
	}
	base, err := r.decoder.Decode(obj)
	if err != nil {
		return err
	}
	return r.handle(base)
}

func (r *FiberReplyHandler) handle(base resppkg.Response) error {
	// first transform the response
	data, err := r.ResponseTransformOpt.Transform(base)
	if err != nil {
		return err
	}
	result, err := r.ResponseMarshalOpt.Marshal(data)
	if err != nil {
		return err
	}
	r.ctx.Set("Content-Type", base.MIMEType)
	return r.ctx.
		Status(base.Code).
		Send(result) //send raw bytes
}
