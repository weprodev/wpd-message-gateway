package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

const HeaderInternalSecret = "X-Internal-Secret"

// InternalIngestSecret requires MESSAGE_INTERNAL_INGEST_SECRET when set; otherwise allows (local dev).
func InternalIngestSecret() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			want := strings.TrimSpace(os.Getenv("MESSAGE_INTERNAL_INGEST_SECRET"))
			if want == "" {
				return next(c)
			}
			if c.Request().Header.Get(HeaderInternalSecret) != want {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or missing internal secret")
			}
			return next(c)
		}
	}
}
