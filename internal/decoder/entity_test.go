package decoder

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	errpkg "github.com/gobeetle/reply/internal/err"
)

func TestDecodePreservesStatusThroughWrap(t *testing.T) {
	inner := errpkg.New(errors.New("bad input")).WithCode(http.StatusBadRequest)
	wrapped := fmt.Errorf("reroute failed: %w", inner)

	resp, err := NewDefaultDecoder().Decode(wrapped)
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}
	if len(resp.Errors) != 1 || resp.Errors[0].Error() != "bad input" {
		t.Fatalf("expected inner error preserved, got %#v", resp.Errors)
	}
}

func TestDecodeDirectErrorReply(t *testing.T) {
	direct := errpkg.New(errors.New("missing")).WithCode(http.StatusNotFound)

	resp, err := NewDefaultDecoder().Decode(direct)
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}
	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.Code)
	}
}

func TestDecodePlainErrorDefaultsBadRequest(t *testing.T) {
	resp, err := NewDefaultDecoder().Decode(errors.New("boom"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}
}
