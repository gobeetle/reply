package fiber

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	errpkg "github.com/gobeetle/reply/internal/err"
	"github.com/gobeetle/reply/internal/strip"
	"github.com/gofiber/fiber/v2"
)

func TestErrorHandlerDefaultStrips5xx(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(ErrorHandlerConfig{}),
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return errpkg.ServiceFailed(errors.New("pq: secret db detail"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status: want %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("json: %v body=%s", err, body)
	}
	if _, ok := payload["error"]; ok {
		t.Fatalf("expected error field stripped, body=%s", body)
	}
	if msg, _ := payload["message"].(string); msg != "service failed" {
		t.Fatalf("expected message preserved, body=%s", body)
	}
}

func TestErrorHandlerDisableStrip(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(ErrorHandlerConfig{
			Strip: strip.Config{Disable: true},
		}),
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return errpkg.ServiceFailed(errors.New("pq: secret db detail"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("json: %v body=%s", err, body)
	}
	errs, ok := payload["error"].([]any)
	if !ok || len(errs) == 0 {
		t.Fatalf("expected error details present when Strip.Disable, body=%s", body)
	}
}

func TestErrorHandlerCustomFallbackMessage(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(ErrorHandlerConfig{
			Strip: strip.Config{FallbackMessage: "something went wrong"},
		}),
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return errpkg.New(errors.New("secret")).WithCode(http.StatusInternalServerError)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("json: %v body=%s", err, body)
	}
	if msg, _ := payload["message"].(string); msg != "something went wrong" {
		t.Fatalf("expected custom fallback, body=%s", body)
	}
}

func TestErrorHandlerHookSeesOriginalError(t *testing.T) {
	var seen string
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(ErrorHandlerConfig{
			Hook: func(c *fiber.Ctx, err error) {
				seen = err.Error()
			},
		}),
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return errpkg.ServiceFailed(errors.New("pq: secret db detail"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()
	if seen != "pq: secret db detail" {
		t.Fatalf("hook should see original cause, got %q", seen)
	}
	body, _ := io.ReadAll(resp.Body)
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("json: %v body=%s", err, body)
	}
	if _, ok := payload["error"]; ok {
		t.Fatalf("response should still strip details, body=%s", body)
	}
}

// Nested call stack used by TestNestedWrapHookAndStrip:
//   handler -> usecase -> service -> db
func dbQuery() error {
	return errors.New("pq: relation users does not exist")
}

func serviceFetchUser() error {
	if err := dbQuery(); err != nil {
		return errpkg.ServiceFailed(err)
	}
	return nil
}

func usecaseGetUser() error {
	if err := serviceFetchUser(); err != nil {
		return fmt.Errorf("usecase get user: %w", err)
	}
	return nil
}

func handlerGetUser() error {
	if err := usecaseGetUser(); err != nil {
		return fmt.Errorf("handler get user: %w", err)
	}
	return nil
}

func TestNestedWrapHookAndStrip(t *testing.T) {
	var (
		hookLog string
		asCode  int
		asMsg   []string
		found   bool
	)

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(ErrorHandlerConfig{
			Hook: func(c *fiber.Ctx, err error) {
				hookLog = err.Error()
				var er *errpkg.ErrorReply
				if errors.As(err, &er) {
					found = true
					asCode = er.StatusCode()
					asMsg = er.Message()
				}
			},
		}),
	})
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return handlerGetUser()
	})

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	// Hook sees the full nested chain for logging/tracing.
	wantLog := "handler get user: usecase get user: pq: relation users does not exist"
	if hookLog != wantLog {
		t.Fatalf("hook log:\n  want: %q\n  got:  %q", wantLog, hookLog)
	}
	if !found {
		t.Fatal("hook errors.As should find ErrorReply through %w chain")
	}
	if asCode != http.StatusInternalServerError {
		t.Fatalf("ErrorReply status: want %d, got %d", http.StatusInternalServerError, asCode)
	}
	if len(asMsg) != 1 || asMsg[0] != "service failed" {
		t.Fatalf("ErrorReply message: want [service failed], got %#v", asMsg)
	}

	// Client JSON keeps status/message but strips internal cause.
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("http status: want %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("json: %v body=%s", err, body)
	}
	if _, ok := payload["error"]; ok {
		t.Fatalf("client must not see db detail, body=%s", body)
	}
	if strings.Contains(string(body), "pq:") || strings.Contains(string(body), "relation users") {
		t.Fatalf("client body leaked secret, body=%s", body)
	}
	if msg, _ := payload["message"].(string); msg != "service failed" {
		t.Fatalf("client message: want service failed, body=%s", body)
	}
}

func TestNestedWrapInvalidRequestPreservesClientDetail(t *testing.T) {
	parseInput := func() error {
		return errors.New("field email is required")
	}
	usecase := func() error {
		if err := parseInput(); err != nil {
			return errpkg.InvalidRequest(err)
		}
		return nil
	}
	handler := func() error {
		if err := usecase(); err != nil {
			return fmt.Errorf("create user: %w", err)
		}
		return nil
	}

	var hookLog string
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(ErrorHandlerConfig{
			Hook: func(c *fiber.Ctx, err error) {
				hookLog = err.Error()
			},
		}),
	})
	app.Post("/users", func(c *fiber.Ctx) error {
		return handler()
	})

	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if hookLog != "create user: field email is required" {
		t.Fatalf("hook log: got %q", hookLog)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("http status: want %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("json: %v body=%s", err, body)
	}
	// 4xx is not stripped — client may see the validation detail.
	errs, ok := payload["error"].([]any)
	if !ok || len(errs) == 0 {
		t.Fatalf("expected client-facing 4xx errors, body=%s", body)
	}
	if errs[0] != "field email is required" {
		t.Fatalf("expected validation detail, body=%s", body)
	}
}
