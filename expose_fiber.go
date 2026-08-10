package reply

import (
	"github.com/gobeetle/reply/internal/strip"
	fiberpkg "github.com/gobeetle/reply/pkg/handler/fiber"
	"github.com/gofiber/fiber/v2"
)

type (
	FiberReplyHandler       = fiberpkg.FiberReplyHandler
	FiberErrorHandlerConfig = fiberpkg.ErrorHandlerConfig
	FiberErrorHook          = fiberpkg.ErrorHook
	ErrorStripper           = strip.ErrorStripper
	StripConfig             = strip.Config
)

var (
	StripServerErrorDetails     = strip.StripServerErrorDetails
	ServerErrorDetails          = strip.ServerErrorDetails
	NoStrip                     = strip.NoStrip
	DefaultStripFallbackMessage = strip.DefaultFallbackMessage
)

// NewFiber creates a new fiber reply handler
func NewFiber(c *fiber.Ctx) *FiberReplyHandler {
	return fiberpkg.NewFiberHandler(c)
}

// FiberErrorHandler creates a Fiber error handler.
// Default: hide internal details on 5xx responses.
//
//	reply.FiberErrorHandler(reply.FiberErrorHandlerConfig{
//	    Hook: func(c *fiber.Ctx, err error) {
//	        logger.ErrorLog(c.Context(), err, nil)
//	    },
//	    Transform: reply.TransformConfig{
//	        Transformer: func(source reply.DefaultResponse) (reply.ErrorCoder, error) {
//	            return &source, nil
//	        },
//	    },
//	    Strip: reply.StripConfig{
//	        FallbackMessage: "something went wrong",
//	    },
//	})
func FiberErrorHandler(cfgs ...FiberErrorHandlerConfig) fiber.ErrorHandler {
	var cfg FiberErrorHandlerConfig
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}
	return fiberpkg.ErrorHandler(cfg)
}
