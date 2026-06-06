package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

const HeaderInternalSecret = "X-Internal-Secret"

// InternalIngestSecret protects internal ingest routes.
// Production must set MESSAGE_INTERNAL_INGEST_SECRET.
// Local automation may set MESSAGE_INTERNAL_INGEST_ALLOW_OPEN=true when no secret is configured.
func InternalIngestSecret() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			want := strings.TrimSpace(os.Getenv("MESSAGE_INTERNAL_INGEST_SECRET"))
			if want == "" {
				if strings.EqualFold(strings.TrimSpace(os.Getenv("MESSAGE_INTERNAL_INGEST_ALLOW_OPEN")), "true") {
					return next(c)
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "internal ingest not configured")
			}
			if c.Request().Header.Get(HeaderInternalSecret) != want {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or missing internal secret")
			}
			return next(c)
		}
	}
}
