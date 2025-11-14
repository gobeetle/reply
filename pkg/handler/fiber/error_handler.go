package fiber

import (
	"github.com/gobeetle/reply/internal/transform"
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler creates a generic error handler for fiber application
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if err == nil {
			return nil
		}
		return NewFiberHandler(c).JSON(err)
	}
}

// ErrorHandlerWithTransform creates an error handler with a custom transform function
func ErrorHandlerWithTransform(opt transform.ResponseTransformOpt) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if err == nil {
			return nil
		}
		return NewFiberHandler(c).
			WithResponseTransformOpt(opt).
			JSON(err)
	}
}
