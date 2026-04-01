package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/authjwt"
)

const (
	// PortalUserIDKey is the context key for JWT-authenticated portal user ID.
	PortalUserIDKey contextKey = "portal_user_id"
	// PortalUserEmailKey carries the email claim for convenience.
	PortalUserEmailKey contextKey = "portal_user_email"
)

// PortalJWT validates Bearer JWTs issued for the portal (not API client credentials).
func PortalJWT(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if secret == "" {
				return echo.NewHTTPError(http.StatusServiceUnavailable, "portal JWT not configured")
			}
			h := strings.TrimSpace(c.Request().Header.Get(echo.HeaderAuthorization))
			const prefix = "Bearer "
			if len(h) < len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing bearer token")
			}
			raw := strings.TrimSpace(h[len(prefix):])
			claims, err := authjwt.Parse(raw, secret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, PortalUserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, PortalUserEmailKey, claims.Email)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

// GetPortalUserID returns the user ID from PortalJWT middleware.
func GetPortalUserID(ctx context.Context) string {
	id, _ := ctx.Value(PortalUserIDKey).(string)
	return id
}

// GetPortalUserEmail returns the email claim from PortalJWT middleware.
func GetPortalUserEmail(ctx context.Context) string {
	email, _ := ctx.Value(PortalUserEmailKey).(string)
	return email
}
