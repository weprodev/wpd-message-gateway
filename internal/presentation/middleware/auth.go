package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/go-pkg/crypto"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

type contextKey string

const (
	// WorkspaceIDKey is the context key for the authenticated workspace ID.
	WorkspaceIDKey contextKey = "workspace_id"
	// APIKeyIDKey is the context key for the authenticated API key row ID.
	APIKeyIDKey contextKey = "api_key_id"
	// APIKeyNameKey is the context key for the API key display name (product/service label).
	APIKeyNameKey contextKey = "api_key_name"
	// HeaderWorkspaceKey is the HTTP header carrying the workspace unique_key (stable slug).
	HeaderWorkspaceKey = "X-Workspace-Key"
)

// APIKeyAuthMiddleware validates client credentials, optional workspace binding via X-Workspace-Key,
// and records last_used_at on the API key.
func APIKeyAuthMiddleware(apiKeyRepo port.APIKeyRepository, workspaceRepo port.WorkspaceRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			clientID, secret, ok := c.Request().BasicAuth()
			if !ok {
				authHeader := c.Request().Header.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					token := strings.TrimPrefix(authHeader, "Bearer ")
					parts := strings.SplitN(token, ":", 2)
					if len(parts) == 2 {
						clientID = parts[0]
						secret = parts[1]
					}
				}
			}

			if clientID == "" || secret == "" {
				slog.WarnContext(c.Request().Context(), "API key auth failed: missing credentials")
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: missing credentials")
			}

			apiKey, err := apiKeyRepo.GetByClientID(c.Request().Context(), clientID)
			if err != nil || !apiKey.IsActive {
				slog.WarnContext(c.Request().Context(), "API key auth failed: invalid or inactive client ID", "client_id", clientID)
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: invalid or inactive client ID")
			}

			if !crypto.CheckSecretHash(secret, apiKey.ClientSecretHash) {
				slog.WarnContext(c.Request().Context(), "API key auth failed: invalid secret", "client_id", clientID)
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: invalid secret")
			}

			wsKey := strings.TrimSpace(c.Request().Header.Get(HeaderWorkspaceKey))
			if wsKey == "" {
				slog.WarnContext(c.Request().Context(), "API key auth failed: missing workspace key header", "client_id", clientID)
				return echo.NewHTTPError(http.StatusBadRequest, "Missing "+HeaderWorkspaceKey+" (workspace unique_key)")
			}
			ws, werr := workspaceRepo.GetByUniqueKey(c.Request().Context(), wsKey)
			if werr != nil || ws == nil {
				slog.WarnContext(c.Request().Context(), "API key auth failed: unknown workspace key", "client_id", clientID, "workspace_key", wsKey)
				return echo.NewHTTPError(http.StatusForbidden, "Unknown workspace key")
			}
			if ws.ID != apiKey.WorkspaceID {
				slog.WarnContext(c.Request().Context(), "API key auth failed: workspace mismatch", "client_id", clientID, "apiKey_workspace_id", apiKey.WorkspaceID, "requested_workspace_id", ws.ID)
				return echo.NewHTTPError(http.StatusForbidden, "API key is not valid for this workspace")
			}

			if err := apiKeyRepo.UpdateLastUsedAt(c.Request().Context(), apiKey.ID); err != nil {
				// Non-fatal: continue request
				_ = err
			}

			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, WorkspaceIDKey, apiKey.WorkspaceID)
			ctx = context.WithValue(ctx, APIKeyIDKey, apiKey.ID)
			ctx = context.WithValue(ctx, APIKeyNameKey, apiKey.Name)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

// GetWorkspaceID extracts the workspace ID injected by APIKeyAuthMiddleware.
func GetWorkspaceID(ctx context.Context) string {
	id, _ := ctx.Value(WorkspaceIDKey).(string)
	return id
}

// GetAPIKeyID extracts the API key ID injected by APIKeyAuthMiddleware.
func GetAPIKeyID(ctx context.Context) string {
	id, _ := ctx.Value(APIKeyIDKey).(string)
	return id
}

// GetAPIKeyName extracts the API key name (product/service label).
func GetAPIKeyName(ctx context.Context) string {
	name, _ := ctx.Value(APIKeyNameKey).(string)
	return name
}
