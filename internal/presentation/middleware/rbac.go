package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	gogate "github.com/weprodev/wpd-gogate"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

// RequirePermission wraps wpd-gogate's RequirePermission but intercepts checks for public workspaces.
func RequirePermission(gate *gogate.Gate, workspaceRepo port.WorkspaceRepository, permissionName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid := GetPortalUserID(c.Request().Context())
			if uid == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: missing user context")
			}

			wid := c.Param("wid")
			if wid == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Bad Request: missing workspace id")
			}

			// Get the workspace to check its visibility
			ws, err := workspaceRepo.GetByID(c.Request().Context(), wid)
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound, "Workspace not found")
			}

			// If workspace is public and it is a read-only request, allow access
			if ws.Visibility == "public" && strings.HasSuffix(permissionName, ".read") {
				return next(c)
			}

			// Otherwise, proceed with the standard permission check via gogate
			allowed, err := gate.Check(c.Request().Context(), "users", uid, permissionName, wid)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error: authorization check failed")
			}

			if !allowed {
				return echo.NewHTTPError(http.StatusForbidden, "Forbidden: missing permission "+permissionName)
			}

			return next(c)
		}
	}
}
