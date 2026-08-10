package strip

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gobeetle/reply/internal/response"
	"github.com/gobeetle/reply/internal/status"
)

func TestStripServerErrorDetailsClears5xxErrors(t *testing.T) {
	in := response.Response{
		Code:   http.StatusInternalServerError,
		Status: status.New(http.StatusInternalServerError),
		Errors: []error{errors.New("pq: relation users does not exist")},
		Msg:    []string{"service failed"},
	}
	out := StripServerErrorDetails(in)
	if len(out.Errors) != 0 {
		t.Fatalf("expected errors cleared, got %#v", out.Errors)
	}
	if len(out.Msg) != 1 || out.Msg[0] != "service failed" {
		t.Fatalf("expected message preserved, got %#v", out.Msg)
	}
	if out.Code != http.StatusInternalServerError {
		t.Fatalf("expected status preserved, got %d", out.Code)
	}
}

func TestStripServerErrorDetailsDefaultMessage(t *testing.T) {
	in := response.Response{
		Code:   http.StatusBadGateway,
		Errors: []error{errors.New("upstream timeout")},
	}
	out := StripServerErrorDetails(in)
	if len(out.Errors) != 0 {
		t.Fatalf("expected errors cleared, got %#v", out.Errors)
	}
	if len(out.Msg) != 1 || out.Msg[0] != DefaultFallbackMessage {
		t.Fatalf("expected default message, got %#v", out.Msg)
	}
}

func TestServerErrorDetailsCustomFallback(t *testing.T) {
	in := response.Response{
		Code:   http.StatusInternalServerError,
		Errors: []error{errors.New("secret")},
	}
	out := ServerErrorDetails("something went wrong")(in)
	if len(out.Msg) != 1 || out.Msg[0] != "something went wrong" {
		t.Fatalf("expected custom fallback, got %#v", out.Msg)
	}
}

func TestConfigResolveDisable(t *testing.T) {
	if (Config{Disable: true}).Resolve() != nil {
		t.Fatal("expected nil stripper when Disable is true")
	}
}

func TestConfigResolveCustomFallback(t *testing.T) {
	stripper := (Config{FallbackMessage: "unavailable"}).Resolve()
	if stripper == nil {
		t.Fatal("expected stripper")
	}
	out := stripper(response.Response{
		Code:   http.StatusInternalServerError,
		Errors: []error{errors.New("secret")},
	})
	if len(out.Msg) != 1 || out.Msg[0] != "unavailable" {
		t.Fatalf("expected custom fallback from config, got %#v", out.Msg)
	}
}

func TestStripServerErrorDetailsLeaves4xx(t *testing.T) {
	in := response.Response{
		Code:   http.StatusBadRequest,
		Errors: []error{errors.New("field x required")},
		Msg:    []string{"validation failed"},
	}
	out := StripServerErrorDetails(in)
	if len(out.Errors) != 1 {
		t.Fatalf("expected 4xx errors preserved, got %#v", out.Errors)
	}
}
