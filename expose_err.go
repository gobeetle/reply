package reply

import (
	errpkg "github.com/gobeetle/reply/internal/err"
)

type (
	ErrorReply = errpkg.ErrorReply
)

var (
	NewError       = errpkg.New
	NewErrorString = errpkg.NewString

	InvalidRequest   = errpkg.InvalidRequest
	ValidationFailed = errpkg.ValidationFailed
	Unauthorized     = errpkg.Unauthorized
	Forbidden        = errpkg.Forbidden
	NotFound         = errpkg.NotFound
	Conflict         = errpkg.Conflict
	Internal         = errpkg.Internal
	ServiceFailed    = errpkg.ServiceFailed
)
