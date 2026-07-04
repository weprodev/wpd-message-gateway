package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/internal/presentation/dto"
)

type PortalIntegrationHandler struct {
	svc *service.PortalService
}

func NewPortalIntegrationHandler(svc *service.PortalService) *PortalIntegrationHandler {
	return &PortalIntegrationHandler{svc: svc}
}

func (h *PortalIntegrationHandler) ListIntegrations(c echo.Context) error {
	wid := c.Param("wid")
	list, err := h.svc.ListIntegrations(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list integrations", "error", err)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to list integrations")
	}
	out := make([]dto.Integration, 0, len(list))
	for _, intg := range list {
		out = append(out, dto.IntegrationFromDomain(intg))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *PortalIntegrationHandler) UpsertIntegration(c echo.Context) error {
	wid := c.Param("wid")
	var body dto.UpsertIntegrationRequest
	if err := c.Bind(&body); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to bind integration body", "error", err, "workspace_id", wid)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	intg, err := body.ToDomain(wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to parse integration config JSON", "error", err, "workspace_id", wid, "provider_name", body.ProviderName)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid config JSON")
	}
	if err := h.svc.UpsertIntegration(c.Request().Context(), intg); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to upsert integration", "error", err, "workspace_id", wid, "provider_name", body.ProviderName)
		return safeHTTPError(err, http.StatusBadRequest, "failed to upsert integration")
	}
	return c.JSON(http.StatusOK, dto.IntegrationFromDomain(*intg))
}

func (h *PortalIntegrationHandler) DeleteIntegration(c echo.Context) error {
	wid := c.Param("wid")
	iid := c.Param("iid")
	if err := h.svc.DeleteIntegration(c.Request().Context(), wid, iid); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to delete integration", "error", err, "workspace_id", wid, "integration_id", iid)
		return safeHTTPError(err, http.StatusBadRequest, "failed to delete integration")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalIntegrationHandler) GetProviderConfigFields(c echo.Context) error {
	wid := c.Param("wid")
	name := c.Param("name")
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "provider name required")
	}
	fields, err := h.svc.GetProviderFields(c.Request().Context(), name)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to get provider config fields", "error", err, "workspace_id", wid, "provider_name", name)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to get provider config fields")
	}
	return c.JSON(http.StatusOK, fields)
}

func (h *PortalIntegrationHandler) ListProviders(c echo.Context) error {
	wid := c.Param("wid")
	providers, err := h.svc.ListProviders(c.Request().Context())
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list providers", "error", err, "workspace_id", wid)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to list providers")
	}
	return c.JSON(http.StatusOK, providers)
}
