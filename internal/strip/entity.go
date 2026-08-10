package strip

import (
	"net/http"

	"github.com/gobeetle/reply/internal/response"
)

const DefaultFallbackMessage = "internal server error"

// ErrorStripper mutates a decoded response before transform/marshal.
type ErrorStripper func(source response.Response) response.Response

// Config controls how server error details are redacted for clients.
// Zero value: strip 5xx errors[] and use DefaultFallbackMessage when message is empty.
type Config struct {
	// Disable keeps full error details in the JSON response.
	Disable bool

	// Stripper is an optional custom redaction function.
	// Leave nil to use the default 5xx strip.
	Stripper ErrorStripper

	// FallbackMessage is used for 5xx responses that have no message after stripping.
	// Leave empty for DefaultFallbackMessage ("internal server error").
	FallbackMessage string
}

// Resolve returns the stripper to apply, or nil when stripping is disabled.
func (c Config) Resolve() ErrorStripper {
	if c.Disable {
		return nil
	}
	if c.Stripper != nil {
		return c.Stripper
	}
	return ServerErrorDetails(c.FallbackMessage)
}

// ServerErrorDetails clears errors[] on 5xx responses.
// fallback is used only when the response has no message; empty fallback uses DefaultFallbackMessage.
func ServerErrorDetails(fallback string) ErrorStripper {
	if fallback == "" {
		fallback = DefaultFallbackMessage
	}
	return func(source response.Response) response.Response {
		if source.Code < http.StatusInternalServerError {
			return source
		}
		source.Errors = nil
		if len(source.Msg) == 0 {
			source.Msg = []string{fallback}
		}
		return source
	}
}

// StripServerErrorDetails is the default 5xx stripper.
var StripServerErrorDetails = ServerErrorDetails(DefaultFallbackMessage)

// NoStrip leaves the response unchanged.
func NoStrip(source response.Response) response.Response {
	return source
}
