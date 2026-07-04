package presentation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/weprodev/wpd-message-gateway/internal/core/logctx"
)

func TestEnrichRequestIDContext_fromEchoMiddleware(t *testing.T) {
	t.Parallel()

	e := echo.New()
	e.Use(echomiddleware.RequestID())
	e.Use(enrichRequestIDContext())
	e.GET("/", func(c echo.Context) error {
		if logctx.GetRequestID(c.Request().Context()) == "" {
			t.Fatal("expected request_id in context")
		}
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: got %d", rec.Code)
	}
	if rec.Header().Get(echo.HeaderXRequestID) == "" {
		t.Fatal("expected X-Request-ID response header")
	}
}

func TestEnrichRequestIDContext_honorsInboundHeader(t *testing.T) {
	t.Parallel()

	e := echo.New()
	e.Use(enrichRequestIDContext())
	e.GET("/", func(c echo.Context) error {
		if got := logctx.GetRequestID(c.Request().Context()); got != "req-inbound" {
			t.Fatalf("request_id: got %q want req-inbound", got)
		}
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderXRequestID, "req-inbound")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: got %d", rec.Code)
	}
}
