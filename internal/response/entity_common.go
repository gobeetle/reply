package response

import (
	"github.com/gobeetle/reply/internal/status"
)

type CommonResponse struct {
	MIMEType string        `json:"mime_type,omitempty"`
	Code     int           `json:"code,omitempty"`
	Status   status.Status `json:"status,omitempty"`
	Errors   []string      `json:"error,omitempty"`
	Msg      *string       `json:"message,omitempty"`
	Data     any           `json:"data,omitempty"`
}

type CommonDataResponse struct {
	MIMEType string        `json:"mime_type,omitempty"`
	Code     int           `json:"code,omitempty"`
	Status   status.Status `json:"status,omitempty"`
	Msg      *string       `json:"message,omitempty"`
	Data     any           `json:"data" validate:"required"`
}

type CommonErrorResponse struct {
	MIMEType string        `json:"mime_type,omitempty"`
	Code     int           `json:"code,omitempty"`
	Status   status.Status `json:"status,omitempty"`
	Errors   []string      `json:"error" validate:"required"`
	Msg      *string       `json:"message,omitempty"`
}

type CommonMessageResponse struct {
	MIMEType string        `json:"mime_type,omitempty"`
	Code     int           `json:"code,omitempty"`
	Status   status.Status `json:"status,omitempty"`
	Msg      *string       `json:"message" validate:"required"`
}
