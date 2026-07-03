package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
)

// PortalHandler exposes REST endpoints for the web portal (JWT, not API keys).
type PortalHandler struct {
	svc     *service.PortalService
	gateway *service.GatewayService
	logs    port.MessageRequestLogRepository
	helper  *SendHelper
}

func NewPortalHandler(
	svc *service.PortalService,
	gateway *service.GatewayService,
	logs port.MessageRequestLogRepository,
) *PortalHandler {
	return &PortalHandler{
		svc:     svc,
		gateway: gateway,
		logs:    logs,
		helper:  NewSendHelper(gateway),
	}
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

type apiKeyPublic struct {
	ID          string     `json:"id"`
	WorkspaceID string     `json:"workspace_id"`
	ClientID    string     `json:"client_id"`
	Name        string     `json:"name"`
	IsActive    bool       `json:"is_active"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type createAPIKeyBody struct {
	Name string `json:"name"`
}

func toAPIKeyPublic(k domain.APIKey) apiKeyPublic {
	return apiKeyPublic{
		ID: k.ID, WorkspaceID: k.WorkspaceID, ClientID: k.ClientID, Name: k.Name,
		IsActive: k.IsActive, LastUsedAt: k.LastUsedAt, CreatedAt: k.CreatedAt, ExpiresAt: k.ExpiresAt,
	}
}

func (h *PortalHandler) ListAPIKeys(c echo.Context) error {
	wid := c.Param("wid")
	keys, err := h.svc.ListAPIKeys(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list api keys", "error", err)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to list api keys")
	}
	out := make([]apiKeyPublic, 0, len(keys))
	for _, k := range keys {
		out = append(out, toAPIKeyPublic(k))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *PortalHandler) CreateAPIKey(c echo.Context) error {
	wid := c.Param("wid")
	var body createAPIKeyBody
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
	pub := toAPIKeyPublic(*k)
	return c.JSON(http.StatusCreated, map[string]any{
		"id": pub.ID, "workspace_id": pub.WorkspaceID, "client_id": pub.ClientID,
		"name": pub.Name, "is_active": pub.IsActive, "last_used_at": pub.LastUsedAt,
		"created_at": pub.CreatedAt, "expires_at": pub.ExpiresAt,
		"client_secret": secret,
	})
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
	return c.JSON(http.StatusOK, map[string]string{"client_secret": secret})
}

type messageRequestLogDTO struct {
	ID           string    `json:"id"`
	WorkspaceID  string    `json:"workspace_id"`
	APIKeyID     string    `json:"api_key_id,omitempty"`
	ChannelType  string    `json:"channel_type"`
	HTTPMethod   string    `json:"http_method"`
	StatusCode   int       `json:"status_code"`
	Endpoint     string    `json:"endpoint"`
	ProviderName string    `json:"provider_name,omitempty"`
	RequestID    string    `json:"request_id,omitempty"`
	DurationMs   int       `json:"duration_ms,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	SourceName   string    `json:"source_name"`
	ClientID     string    `json:"client_id,omitempty"`
}

func toMessageRequestLogDTO(l domain.MessageRequestLogWithSource) messageRequestLogDTO {
	return messageRequestLogDTO{
		ID:           l.ID,
		WorkspaceID:  l.WorkspaceID,
		APIKeyID:     l.APIKeyID,
		ChannelType:  l.ChannelType,
		HTTPMethod:   l.HTTPMethod,
		StatusCode:   l.StatusCode,
		Endpoint:     l.Endpoint,
		ProviderName: l.ProviderName,
		RequestID:    l.RequestID,
		DurationMs:   l.DurationMs,
		ErrorMessage: l.ErrorMessage,
		CreatedAt:    l.CreatedAt,
		SourceName:   l.SourceName,
		ClientID:     l.ClientID,
	}
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
	out := make([]messageRequestLogDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, toMessageRequestLogDTO(r))
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out, "total": total})
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
	var body map[string]string
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
