package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/go-pkg/crypto"

	"github.com/weprodev/wpd-message-gateway/internal/core/authjwt"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

const (
	// HeaderWorkspaceAPIClientID is the client_id of a workspace API key (used with JWT for portal inbox).
	HeaderWorkspaceAPIClientID = "X-Api-Client-Id"
	// HeaderWorkspaceAPISecret is the plaintext secret for that API key.
	HeaderWorkspaceAPISecret = "X-Api-Client-Secret"
	// QueryAccessToken is an optional JWT query param for endpoints that cannot set headers (e.g. EventSource); prefer Authorization header.
	QueryAccessToken = "access_token"
	// QueryClientID / QueryClientSecret mirror headers for SSE clients that cannot send custom headers.
	QueryClientID     = "client_id"
	QueryClientSecret = "client_secret"
)

// PortalJWTBearerOrQuery validates the portal JWT from Authorization: Bearer <token>,
// or from query access_token on GET requests (limited compatibility for SSE clients).
func PortalJWTBearerOrQuery(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if secret == "" {
				slog.ErrorContext(c.Request().Context(), "Portal JWT auth failed: secret not configured")
				return echo.NewHTTPError(http.StatusServiceUnavailable, "portal JWT not configured")
			}
			raw := bearerTokenFromRequest(c)
			if raw == "" {
				slog.WarnContext(c.Request().Context(), "Portal JWT auth failed: missing bearer token")
				return echo.NewHTTPError(http.StatusUnauthorized, "missing bearer token")
			}
			claims, err := authjwt.Parse(raw, secret)
			if err != nil {
				slog.WarnContext(c.Request().Context(), "Portal JWT auth failed: invalid token", "error", err)
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

func bearerTokenFromRequest(c echo.Context) string {
	h := strings.TrimSpace(c.Request().Header.Get(echo.HeaderAuthorization))
	const prefix = "Bearer "
	if len(h) >= len(prefix) && strings.EqualFold(h[:len(prefix)], prefix) {
		return strings.TrimSpace(h[len(prefix):])
	}
	if c.Request().Method == http.MethodGet {
		if q := strings.TrimSpace(c.QueryParam(QueryAccessToken)); q != "" {
			return q
		}
	}
	return ""
}

// RequireWorkspaceMember ensures the JWT user is a member of :wid (must run after PortalJWTBearerOrQuery).
func RequireWorkspaceMember(members port.WorkspaceMemberRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			wid := c.Param("wid")
			if wid == "" {
				slog.WarnContext(c.Request().Context(), "RequireWorkspaceMember failed: missing workspace id param")
				return echo.NewHTTPError(http.StatusBadRequest, "missing workspace id")
			}
			uid := GetPortalUserID(c.Request().Context())
			if uid == "" {
				slog.WarnContext(c.Request().Context(), "RequireWorkspaceMember failed: missing user context", "workspace_id", wid)
				return echo.NewHTTPError(http.StatusUnauthorized, "missing user context")
			}
			if _, err := members.GetRole(c.Request().Context(), wid, uid); err != nil {
				slog.WarnContext(c.Request().Context(), "RequireWorkspaceMember failed: user is not a member", "workspace_id", wid, "user_id", uid, "error", err)
				return echo.NewHTTPError(http.StatusForbidden, "not a member of this workspace")
			}
			return next(c)
		}
	}
}

// RequireWorkspaceAPIKey validates X-Api-Client-Id + X-Api-Client-Secret against the workspace API key table.
// The key must belong to :wid. Headers are preferred; GET may use client_id and client_secret query params for SSE.
func RequireWorkspaceAPIKey(apiKeys port.APIKeyRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			wid := c.Param("wid")
			if wid == "" {
				slog.WarnContext(c.Request().Context(), "RequireWorkspaceAPIKey failed: missing workspace id param")
				return echo.NewHTTPError(http.StatusBadRequest, "missing workspace id")
			}
			clientID := strings.TrimSpace(c.Request().Header.Get(HeaderWorkspaceAPIClientID))
			secret := c.Request().Header.Get(HeaderWorkspaceAPISecret)
			if clientID == "" {
				clientID = strings.TrimSpace(c.QueryParam(QueryClientID))
			}
			if secret == "" {
				secret = c.QueryParam(QueryClientSecret)
			}
			if clientID == "" || secret == "" {
				slog.WarnContext(c.Request().Context(), "RequireWorkspaceAPIKey failed: missing credentials", "workspace_id", wid)
				return echo.NewHTTPError(http.StatusUnauthorized, "missing API key credentials (X-Api-Client-Id and X-Api-Client-Secret)")
			}
			key, err := apiKeys.GetByClientID(c.Request().Context(), clientID)
			if err != nil || !key.IsActive {
				slog.WarnContext(c.Request().Context(), "RequireWorkspaceAPIKey failed: invalid or inactive key", "workspace_id", wid, "client_id", clientID)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
			}
			if !crypto.CheckSecretHash(secret, key.ClientSecretHash) {
				slog.WarnContext(c.Request().Context(), "RequireWorkspaceAPIKey failed: invalid secret", "workspace_id", wid, "client_id", clientID)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key secret")
			}
			if key.WorkspaceID != wid {
				slog.WarnContext(c.Request().Context(), "RequireWorkspaceAPIKey failed: workspace mismatch", "workspace_id", wid, "client_id", clientID, "apiKey_workspace_id", key.WorkspaceID)
				return echo.NewHTTPError(http.StatusForbidden, "API key does not belong to this workspace")
			}
			_ = apiKeys.UpdateLastUsedAt(c.Request().Context(), key.ID)
			return next(c)
		}
	}
}
