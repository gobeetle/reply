package response

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/gobeetle/reply/internal/status"
	"github.com/gobeetle/reply/internal/utils"
)

type Response struct {
	MIMEType string
	Code     int
	Status   status.Status
	Errors   []error
	Msg      []string
	Data     any
}

func (r *Response) Error() string {
	if len(r.Errors) == 0 {
		return ""
	}
	return errors.Join(r.Errors...).Error()
}

func (r *Response) StatusCode() int {
	return r.Code
}

func (r *Response) MarshalJSON() ([]byte, error) {
	msg := ""
	if len(r.Msg) > 0 {
		msg = strings.Join(r.Msg, ", ")
	}
	return json.Marshal(CommonResponse{
		MIMEType: r.MIMEType,
		Code:     r.Code,
		Status:   r.Status,
		Errors:   utils.ErrorsToStrings(r.Errors),
		Msg:      &msg,
		Data:     r.Data,
	})
}

// var _ iface.ErrorCoder = (*DefaultResponse)(nil)
// var _ iface.StatusCoder = (*DefaultResponse)(nil)
// var _ error = (*DefaultResponse)(nil)
