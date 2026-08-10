package err

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestUseCaseConstructors(t *testing.T) {
	cause := errors.New("detail")

	cases := []struct {
		name    string
		got     *ErrorReply
		code    int
		message string
	}{
		{"InvalidRequest", InvalidRequest(cause), http.StatusBadRequest, "invalid request"},
		{"ValidationFailed", ValidationFailed(cause), http.StatusBadRequest, "validation failed"},
		{"Unauthorized", Unauthorized(cause), http.StatusUnauthorized, "unauthorized"},
		{"Forbidden", Forbidden(cause), http.StatusForbidden, "forbidden"},
		{"NotFound", NotFound(cause), http.StatusNotFound, "not found"},
		{"Conflict", Conflict(cause), http.StatusConflict, "conflict"},
		{"Internal", Internal(cause), http.StatusInternalServerError, "internal server error"},
		{"ServiceFailed", ServiceFailed(cause), http.StatusInternalServerError, "service failed"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got.StatusCode() != tc.code {
				t.Fatalf("status: want %d, got %d", tc.code, tc.got.StatusCode())
			}
			msgs := tc.got.Message()
			if len(msgs) != 1 || msgs[0] != tc.message {
				t.Fatalf("message: want %q, got %#v", tc.message, msgs)
			}
			if tc.got.Error() != cause.Error() {
				t.Fatalf("error: want %q, got %q", cause.Error(), tc.got.Error())
			}
		})
	}
}

func TestUseCaseConstructorWrapRoundTrip(t *testing.T) {
	inner := InvalidRequest(errors.New("bad input"))
	wrapped := fmt.Errorf("handler failed: %w", inner)

	var target *ErrorReply
	if !errors.As(wrapped, &target) {
		t.Fatal("expected errors.As to find ErrorReply")
	}
	if target.StatusCode() != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d", http.StatusBadRequest, target.StatusCode())
	}
}
