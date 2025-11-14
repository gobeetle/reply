package decoder

import (
	"net/http"

	"github.com/gobeetle/reply/internal/constant"
	"github.com/gobeetle/reply/internal/iface"
	resppkg "github.com/gobeetle/reply/internal/response"
	"github.com/gobeetle/reply/internal/status"

	datapkg "github.com/gobeetle/reply/internal/data"
	errpkg "github.com/gobeetle/reply/internal/err"
)

type DefaultDecoder struct{}

func NewDefaultDecoder() *DefaultDecoder {
	return &DefaultDecoder{}
}

func (r *DefaultDecoder) Empty() (resppkg.Response, error) {
	return r.Decode(nil)
}

func (r *DefaultDecoder) Decode(obj any) (resppkg.Response, error) {
	if obj == nil {
		// empty response, StatusNoContent
		return r.decode(nil)
	} else if coder, ok := obj.(iface.StatusCoder); ok {
		// if obj is StatusCoder, then decode it
		return r.decode(coder)
	} else if err, ok := obj.(iface.ErrorProvider); ok {
		// if it is just a generic error, then wrap it as ErrorCoder
		return r.decode(errpkg.New(err))
	}
	// for anything else, wrap it as DataCoder
	return r.decode(datapkg.New(obj))
}

func (r *DefaultDecoder) decode(coder iface.StatusCoder) (resppkg.Response, error) {
	base := resppkg.Response{}
	base.MIMEType = constant.MediaTypeJson
	base.Code = 0
	if coder != nil {
		base.Code = coder.StatusCode()
		if d, ok := coder.(iface.DataCoder); ok {
			base.Data = d.Data()
		}
		if e, ok := coder.(iface.ErrorCoder); ok {
			base.Errors = []error{e}
			if errs, ok := e.(iface.Errors); ok {
				errs := errs.Unwrap()
				if len(errs) > 0 {
					base.Errors = errs
				}
			}
			base.MIMEType = constant.MediaTypeProblemJson
		}
		if m, ok := coder.(iface.MessageCoder); ok {
			base.Msg = m.Message()
		}
	}
	// if status code is not set, derive it from the payload
	if base.Code == 0 {
		switch {
		case len(base.Errors) > 0:
			base.Code = http.StatusBadRequest
		case base.Data == nil && base.Msg == nil:
			base.Code = http.StatusNoContent
		default:
			base.Code = http.StatusOK
		}
	}
	base.Status = status.New(base.Code)
	return base, nil
}
