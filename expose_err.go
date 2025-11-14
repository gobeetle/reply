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
)
