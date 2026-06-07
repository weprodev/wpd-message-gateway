package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
)

type PortalIntegrationHandler struct {
	svc *service.PortalService
}

func NewPortalIntegrationHandler(svc *service.PortalService) *PortalIntegrationHandler {
	return &PortalIntegrationHandler{svc: svc}
}

type integrationBody struct {
	ChannelType  string          `json:"channel_type"`
	ProviderName string          `json:"provider_name"`
	Config       json.RawMessage `json:"config"`
	Status       string          `json:"status"`
	IsDefault    bool            `json:"is_default"`
}

func integrationToJSON(intg domain.Integration) map[string]any {
	return map[string]any{
		"id":            intg.ID,
		"workspace_id":  intg.WorkspaceID,
		"channel_type":  intg.ChannelType,
		"provider_name": intg.ProviderName,
		"config":        json.RawMessage(intg.Config),
		"status":        intg.Status,
		"is_default":    intg.IsDefault,
		"created_at":    intg.CreatedAt,
		"updated_at":    intg.UpdatedAt,
	}
}

func (h *PortalIntegrationHandler) ListIntegrations(c echo.Context) error {
	wid := c.Param("wid")
	list, err := h.svc.ListIntegrations(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list integrations", "error", err)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to list integrations")
	}
	out := make([]map[string]any, 0, len(list))
	for _, intg := range list {
		out = append(out, integrationToJSON(intg))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *PortalIntegrationHandler) UpsertIntegration(c echo.Context) error {
	wid := c.Param("wid")
	var body integrationBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	cfg, err := service.IntegrationConfigJSON(body.Config)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid config JSON")
	}
	status := body.Status
	if status == "" {
		status = "connected"
	}
	intg := &domain.Integration{
		WorkspaceID:  wid,
		ChannelType:  body.ChannelType,
		ProviderName: body.ProviderName,
		Config:       cfg,
		Status:       status,
		IsDefault:    body.IsDefault,
	}
	if err := h.svc.UpsertIntegration(c.Request().Context(), intg); err != nil {
		return safeHTTPError(err, http.StatusBadRequest, "failed to upsert integration")
	}
	return c.JSON(http.StatusOK, integrationToJSON(*intg))
}

func (h *PortalIntegrationHandler) DeleteIntegration(c echo.Context) error {
	wid := c.Param("wid")
	iid := c.Param("iid")
	if err := h.svc.DeleteIntegration(c.Request().Context(), wid, iid); err != nil {
		return safeHTTPError(err, http.StatusBadRequest, "failed to delete integration")
	}
	return c.NoContent(http.StatusNoContent)
}
