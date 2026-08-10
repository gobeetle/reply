package err

import "net/http"

// Use-case constructors: intent-named ErrorReply helpers with status + default message.
// Callers can still chain WithMessage / WithError to override.

func InvalidRequest(errs ...error) *ErrorReply {
	return New(errs...).
		WithCode(http.StatusBadRequest).
		WithMessage("invalid request")
}

func ValidationFailed(errs ...error) *ErrorReply {
	return New(errs...).
		WithCode(http.StatusBadRequest).
		WithMessage("validation failed")
}

func Unauthorized(errs ...error) *ErrorReply {
	return New(errs...).
		WithCode(http.StatusUnauthorized).
		WithMessage("unauthorized")
}

func Forbidden(errs ...error) *ErrorReply {
	return New(errs...).
		WithCode(http.StatusForbidden).
		WithMessage("forbidden")
}

func NotFound(errs ...error) *ErrorReply {
	return New(errs...).
		WithCode(http.StatusNotFound).
		WithMessage("not found")
}

func Conflict(errs ...error) *ErrorReply {
	return New(errs...).
		WithCode(http.StatusConflict).
		WithMessage("conflict")
}

func Internal(errs ...error) *ErrorReply {
	return New(errs...).
		WithCode(http.StatusInternalServerError).
		WithMessage("internal server error")
}

// ServiceFailed is a 500 for a known/downstream service failure
// (vs Internal for unexpected faults).
func ServiceFailed(errs ...error) *ErrorReply {
	return New(errs...).
		WithCode(http.StatusInternalServerError).
		WithMessage("service failed")
}
