package fiber

import (
	"github.com/gobeetle/reply/internal/marshal"
	"github.com/gobeetle/reply/internal/strip"
	"github.com/gobeetle/reply/internal/transform"
	"github.com/gofiber/fiber/v2"
)

// ErrorHook runs before the error is encoded into a response.
// Use for project-specific side effects such as logging.
type ErrorHook func(c *fiber.Ctx, err error)

// ErrorHandlerConfig configures FiberErrorHandler.
// Zero value enables default 5xx stripping.
type ErrorHandlerConfig struct {
	// Hook runs first with the original error (before strip/transform/marshal).
	Hook ErrorHook

	// Transform runs after stripping, before the response is written.
	Transform transform.ResponseTransformOpt

	// Marshal replaces the default JSON marshaler when set.
	Marshal marshal.ResponseMarshalOpt

	// Strip controls redaction of internal error details on 5xx responses.
	Strip strip.Config
}

// ErrorHandler creates a Fiber error handler from cfg.
func ErrorHandler(cfg ErrorHandlerConfig) fiber.ErrorHandler {
	stripper := cfg.Strip.Resolve()
	return func(c *fiber.Ctx, err error) error {
		if err == nil {
			return nil
		}
		if cfg.Hook != nil {
			cfg.Hook(c, err)
		}
		h := NewFiberHandler(c).
			WithResponseTransformOpt(cfg.Transform).
			WithResponseMarshalOpt(cfg.Marshal)
		if stripper != nil {
			h = h.WithErrorStripper(stripper)
		}
		return h.JSON(err)
	}
}
