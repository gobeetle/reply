package reply

import (
	"github.com/gobeetle/reply/internal/transform"
	fuegopkg "github.com/gobeetle/reply/pkg/handler/fuego"
)

type (
	FuegoReplyHandler = fuegopkg.FuegoReplyHandler
)

// NewFuego creates a new fuego reply handler
func NewFuego() *FuegoReplyHandler { return fuegopkg.NewFuegoHandler() }

// FuegoErrorHandler creates a generic error handler for fuego application
func FuegoErrorHandler() func(err error) error {
	return fuegopkg.ErrorHandler()
}

// FuegoErrorHandlerWithTransformer creates an error handler with a custom transform function
func FuegoErrorHandlerWithTransformer(opt transform.ResponseTransformOpt) func(err error) error {
	return fuegopkg.ErrorHandlerWithTransformer(opt)
}
