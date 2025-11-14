package reply

import (
	"github.com/gobeetle/reply/internal/transform"
	fiberpkg "github.com/gobeetle/reply/pkg/handler/fiber"
	"github.com/gofiber/fiber/v2"
)

type (
	FiberReplyHandler = fiberpkg.FiberReplyHandler
)

// NewFiber creates a new fiber reply handler
func NewFiber(c *fiber.Ctx) *FiberReplyHandler {
	return fiberpkg.NewFiberHandler(c)
}

// FiberErrorHandler creates a generic error handler for fiber application
func FiberErrorHandler() fiber.ErrorHandler {
	return fiberpkg.ErrorHandler()
}

// FiberErrorHandlerWithTransformer creates an error handler with a custom transform function
func FiberErrorHandlerWithTransformer(opt transform.ResponseTransformOpt) fiber.ErrorHandler {
	return fiberpkg.ErrorHandlerWithTransform(opt)
}
