package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	gogate "github.com/weprodev/wpd-gogate"
)

// RequirePermission wraps wpd-gogate's RequirePermission.
// It extracts the portal user ID and workspace ID (:wid path parameter) from the context.
func RequirePermission(gate *gogate.Gate, permissionName string) echo.MiddlewareFunc {
	opts := gogate.MiddlewareOptions{
		ModelType: "users",
		ExtractModelID: func(c echo.Context) (any, error) {
			uid := GetPortalUserID(c.Request().Context())
			if uid == "" {
				return nil, errors.New("missing user context")
			}
			return uid, nil
		},
		ExtractTeamID: func(c echo.Context) (any, error) {
			wid := c.Param("wid")
			if wid == "" {
				return nil, errors.New("missing workspace id")
			}
			return wid, nil
		},
		OnDenied: func(c echo.Context, permissionName string) error {
			return echo.NewHTTPError(http.StatusForbidden, "Forbidden: missing permission "+permissionName)
		},
		OnError: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error: authorization check failed")
		},
	}
	return gogate.RequirePermission(gate, permissionName, &opts)
}
