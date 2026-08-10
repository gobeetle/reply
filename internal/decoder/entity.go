package decoder

import (
	"errors"
	"net/http"

	"github.com/gobeetle/reply/internal/constant"
	"github.com/gobeetle/reply/internal/iface"
	resppkg "github.com/gobeetle/reply/internal/response"
	"github.com/gobeetle/reply/internal/status"

	datapkg "github.com/gobeetle/reply/internal/data"
	errpkg "github.com/gobeetle/reply/internal/err"
)

// DefaultUnknownErrorStatus is used when an error has no StatusCode.
const DefaultUnknownErrorStatus = http.StatusInternalServerError

type DefaultDecoder struct {
	// unknownErrorStatus is used when StatusCode is 0 and the payload is an error.
	// Zero means DefaultUnknownErrorStatus (500).
	unknownErrorStatus int
}

func NewDefaultDecoder() *DefaultDecoder {
	return &DefaultDecoder{}
}

// WithUnknownErrorStatus sets the HTTP status for errors without an explicit code.
// Pass 0 to restore the package default (500).
func (r *DefaultDecoder) WithUnknownErrorStatus(code int) *DefaultDecoder {
	r.unknownErrorStatus = code
	return r
}

func (r *DefaultDecoder) resolveUnknownErrorStatus() int {
	if r.unknownErrorStatus != 0 {
		return r.unknownErrorStatus
	}
	return DefaultUnknownErrorStatus
}

func (r *DefaultDecoder) Empty() (resppkg.Response, error) {
	return r.Decode(nil)
}

func (r *DefaultDecoder) Decode(obj any) (resppkg.Response, error) {
	if obj == nil {
		// empty response, StatusNoContent
		return r.decode(nil)
	}
	if coder, ok := obj.(iface.StatusCoder); ok {
		// direct StatusCoder (ErrorReply, DataReply, MessageReply, ...)
		return r.decode(coder)
	}
	if err, ok := obj.(error); ok {
		// walk fmt.Errorf("%w") / Unwrap chains for a nested ErrorCoder
		var ec iface.ErrorCoder
		if errors.As(err, &ec) {
			return r.decode(ec)
		}
		// plain error -> wrap as ErrorCoder
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
			base.Code = r.resolveUnknownErrorStatus()
		case base.Data == nil && base.Msg == nil:
			base.Code = http.StatusNoContent
		default:
			base.Code = http.StatusOK
		}
	}
	base.Status = status.New(base.Code)
	return base, nil
}
