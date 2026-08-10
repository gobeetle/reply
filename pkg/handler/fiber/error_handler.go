package fiber

import (
	"net/http"

	"github.com/gobeetle/reply/internal/decoder"
	"github.com/gobeetle/reply/internal/marshal"
	"github.com/gobeetle/reply/internal/strip"
	"github.com/gobeetle/reply/internal/transform"
	"github.com/gofiber/fiber/v2"
)

// ErrorHook runs before the error is encoded into a response.
// Use for project-specific side effects such as logging.
type ErrorHook func(c *fiber.Ctx, err error)

// ErrorHandlerConfig configures FiberErrorHandler.
// Zero value enables default 5xx stripping and treats uncoded errors as 500.
type ErrorHandlerConfig struct {
	// Hook runs first with the original error (before strip/transform/marshal).
	Hook ErrorHook

	// Transform reshapes the outbound envelope after strip, before marshal.
	Transform transform.Config

	// Marshal replaces the default JSON marshaler when set.
	Marshal marshal.ResponseMarshalOpt

	// Strip controls redaction of internal error details on 5xx responses.
	Strip strip.Config

	// UnknownErrorStatus is used when an error has no StatusCode (plain errors.New / fmt.Errorf).
	// Zero value defaults to 500. Set to 400 if you prefer untyped errors as bad requests.
	UnknownErrorStatus int
}

func (cfg ErrorHandlerConfig) resolveUnknownErrorStatus() int {
	if cfg.UnknownErrorStatus != 0 {
		return cfg.UnknownErrorStatus
	}
	return http.StatusInternalServerError
}

// ErrorHandler creates a Fiber error handler from cfg.
func ErrorHandler(cfg ErrorHandlerConfig) fiber.ErrorHandler {
	stripper := cfg.Strip.Resolve()
	transformOpt := cfg.Transform.Resolve()
	dec := decoder.NewDefaultDecoder().WithUnknownErrorStatus(cfg.resolveUnknownErrorStatus())
	return func(c *fiber.Ctx, err error) error {
		if err == nil {
			return nil
		}
		if cfg.Hook != nil {
			cfg.Hook(c, err)
		}
		h := NewFiberHandler(c).
			WithResponseDecoder(dec).
			WithResponseTransformOpt(transformOpt).
			WithResponseMarshalOpt(cfg.Marshal)
		if stripper != nil {
			h = h.WithErrorStripper(stripper)
		}
		return h.JSON(err)
	}
}
