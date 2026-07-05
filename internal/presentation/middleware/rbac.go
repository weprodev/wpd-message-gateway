package middleware

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	gogate "github.com/weprodev/wpd-gogate"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

// RequirePermission wraps wpd-gogate's RequirePermission but intercepts checks for public workspaces.
func RequirePermission(gate *gogate.Gate, workspaceRepo port.WorkspaceRepository, permissionName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid := GetPortalUserID(c.Request().Context())
			if uid == "" {
				slog.WarnContext(c.Request().Context(), "RBAC check failed: missing user context", "permission", permissionName)
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: missing user context")
			}

			wid := c.Param("wid")
			if wid == "" {
				slog.WarnContext(c.Request().Context(), "RBAC check failed: missing workspace id param", "user_id", uid, "permission", permissionName)
				return echo.NewHTTPError(http.StatusBadRequest, "Bad Request: missing workspace id")
			}

			// Get the workspace to check its visibility
			ws, err := workspaceRepo.GetByID(c.Request().Context(), wid)
			if err != nil {
				slog.ErrorContext(c.Request().Context(), "RBAC check failed: workspace not found", "error", err, "user_id", uid, "workspace_id", wid, "permission", permissionName)
				return echo.NewHTTPError(http.StatusNotFound, "Workspace not found")
			}

			// Public workspaces: guest-safe read permissions only (non-members).
			if ws.Visibility == "public" && domain.IsPublicGuestPermission(permissionName) {
				slog.InfoContext(c.Request().Context(), "RBAC check bypassed: public workspace read access", "user_id", uid, "workspace_id", wid, "permission", permissionName)
				return next(c)
			}

			// Otherwise, proceed with the standard permission check via gogate
			allowed, err := gate.Check(c.Request().Context(), "users", uid, permissionName, "", wid)
			if err != nil {
				slog.ErrorContext(c.Request().Context(), "RBAC check failed: gate check error", "error", err, "user_id", uid, "workspace_id", wid, "permission", permissionName)
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error: authorization check failed")
			}

			if !allowed {
				slog.WarnContext(c.Request().Context(), "RBAC check denied", "user_id", uid, "workspace_id", wid, "permission", permissionName)
				return echo.NewHTTPError(http.StatusForbidden, "Forbidden: missing permission "+permissionName)
			}

			return next(c)
		}
	}
}
