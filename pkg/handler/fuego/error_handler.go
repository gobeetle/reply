package fuego

import (
	"github.com/gobeetle/reply/internal/transform"
)

// ErrorHandler creates a generic error handler for fuego application
func ErrorHandler() func(err error) error {
	return func(err error) error {
		if err == nil {
			return nil
		}
		data, decodeErr := NewFuegoHandler().Decode(err)
		if decodeErr != nil {
			return decodeErr
		}
		return data
	}
}

// ErrorHandlerWithTransformer creates an error handler with a custom transform function
func ErrorHandlerWithTransformer(transform transform.ResponseTransformOpt) func(err error) error {
	return func(err error) error {
		if err == nil {
			return nil
		}
		data, decodeErr := NewFuegoHandler().
			WithResponseTransformOpt(transform).
			Decode(err)
		if decodeErr != nil {
			return decodeErr
		}
		return data
	}
}
