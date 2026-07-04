package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/internal/presentation/dto"
)

// PortalHandler exposes REST endpoints for the web portal (JWT, not API keys).
type PortalHandler struct {
	svc *service.PortalService
}

func NewPortalHandler(svc *service.PortalService) *PortalHandler {
	return &PortalHandler{svc: svc}
}

// safeHTTPError returns an HTTPError with a sanitised message.
// For known domain sentinels it uses a curated message; for all other errors
// it returns a generic message to avoid leaking infrastructure details.
func safeHTTPError(err error, fallbackStatus int, fallbackMessage string) *echo.HTTPError {
	switch {
	case errors.Is(err, port.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, "resource not found")
	case errors.Is(err, port.ErrConflict):
		return echo.NewHTTPError(http.StatusConflict, "resource already exists")
	case errors.Is(err, port.ErrUnauthorized):
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	case errors.Is(err, port.ErrInvalidCredentials):
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
	case errors.Is(err, port.ErrInvalidInput):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return echo.NewHTTPError(fallbackStatus, fallbackMessage)
}

func (h *PortalHandler) ListAPIKeys(c echo.Context) error {
	wid := c.Param("wid")
	keys, err := h.svc.ListAPIKeys(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list api keys", "error", err)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to list api keys")
	}
	out := make([]dto.APIKeyPublic, 0, len(keys))
	for _, k := range keys {
		out = append(out, dto.APIKeyPublicFromDomain(k))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *PortalHandler) CreateAPIKey(c echo.Context) error {
	wid := c.Param("wid")
	var body dto.CreateAPIKeyRequest
	if err := c.Bind(&body); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to bind api key body", "error", err, "workspace_id", wid)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if body.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name required")
	}
	k, secret, err := h.svc.CreateAPIKey(c.Request().Context(), wid, body.Name)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to create api key", "error", err, "workspace_id", wid, "name", body.Name)
		return safeHTTPError(err, http.StatusBadRequest, "failed to create api key")
	}
	return c.JSON(http.StatusCreated, dto.CreateAPIKeyResponseFromDomain(*k, secret))
}

func (h *PortalHandler) DeleteAPIKey(c echo.Context) error {
	wid := c.Param("wid")
	keyID := c.Param("keyId")
	if err := h.svc.DeleteAPIKey(c.Request().Context(), wid, keyID); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to delete api key", "error", err, "workspace_id", wid, "api_key_id", keyID)
		return safeHTTPError(err, http.StatusBadRequest, "failed to delete api key")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalHandler) RegenerateAPIKey(c echo.Context) error {
	wid := c.Param("wid")
	keyID := c.Param("keyId")
	secret, err := h.svc.RegenerateAPIKey(c.Request().Context(), wid, keyID)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to regenerate api key", "error", err, "workspace_id", wid, "api_key_id", keyID)
		return safeHTTPError(err, http.StatusBadRequest, "failed to regenerate api key")
	}
	return c.JSON(http.StatusOK, dto.RegenerateAPIKeyResponseFromDomain(secret))
}

func (h *PortalHandler) ListLogs(c echo.Context) error {
	wid := c.Param("wid")
	q := port.MessageLogQuery{WorkspaceID: wid}
	if ch := c.QueryParam("channel"); ch != "" {
		q.ChannelType = ch
	}
	if lim := c.QueryParam("limit"); lim != "" {
		if n, err := strconv.Atoi(lim); err == nil {
			q.Limit = n
		}
	}
	if off := c.QueryParam("offset"); off != "" {
		if n, err := strconv.Atoi(off); err == nil {
			q.Offset = n
		}
	}
	rows, total, err := h.svc.ListLogs(c.Request().Context(), q)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list logs", "error", err)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to list logs")
	}
	return c.JSON(http.StatusOK, dto.MessageRequestLogListFromDomain(rows, total))
}

func (h *PortalHandler) GetSettings(c echo.Context) error {
	wid := c.Param("wid")
	m, err := h.svc.GetSettings(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to get settings", "error", err)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to get settings")
	}
	return c.JSON(http.StatusOK, m)
}

func (h *PortalHandler) PatchSettings(c echo.Context) error {
	wid := c.Param("wid")
	var body dto.PatchSettingsRequest
	if err := c.Bind(&body); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to bind settings body", "error", err, "workspace_id", wid)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if err := h.svc.PatchSettings(c.Request().Context(), wid, body); err != nil {
		if !errors.Is(err, port.ErrInvalidInput) {
			slog.ErrorContext(c.Request().Context(), "failed to patch settings", "error", err, "workspace_id", wid)
		}
		return safeHTTPError(err, http.StatusBadRequest, "failed to patch settings")
	}
	return c.NoContent(http.StatusNoContent)
}
