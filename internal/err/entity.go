package err

import (
	"encoding/json"
	"errors"

	"github.com/gobeetle/reply/internal/envelope"
	"github.com/gobeetle/reply/internal/iface"
)

type ErrorReply struct {
	code    int
	errs    []error
	message []string
}

func New(errs ...error) *ErrorReply {
	return &ErrorReply{
		errs: errs,
	}
}

func NewString(errs ...string) *ErrorReply {
	var results []error
	for _, err := range errs {
		results = append(results, errors.New(err))
	}
	return New(results...)
}

func (e *ErrorReply) WithCode(code int) *ErrorReply {
	e.code = code
	return e
}

func (e *ErrorReply) WithMessage(message ...string) *ErrorReply {
	e.message = message
	return e
}

func (e *ErrorReply) WithError(errs ...error) *ErrorReply {
	e.errs = errs
	return e
}

func (e *ErrorReply) Error() string {
	if len(e.errs) > 0 {
		return errors.Join(e.errs...).Error()
	}
	return ""
}

func (e *ErrorReply) Unwrap() []error {
	return e.errs
}

func (e *ErrorReply) StatusCode() int {
	return e.code
}

func (e *ErrorReply) Message() []string {
	return e.message
}

func (e *ErrorReply) MarshalJSON() ([]byte, error) {
	var errStrings []string
	if e.errs != nil {
		errStrings = make([]string, 0, len(e.errs))
		for _, er := range e.errs {
			if er != nil {
				errStrings = append(errStrings, er.Error())
			}
		}
	}
	return json.Marshal(
		envelope.ReplyEnvelope{
			Code:    e.code,
			Errors:  errStrings,
			Message: e.message,
		},
	)
}

var (
	_ iface.ErrorCoder    = (*ErrorReply)(nil)
	_ iface.MessageCoder  = (*ErrorReply)(nil)
	_ iface.Errors        = (*ErrorReply)(nil)
	_ iface.JsonMarshaler = (*ErrorReply)(nil)
)
