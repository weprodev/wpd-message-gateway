package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/internal/presentation/dto"
)

type PortalTemplateHandler struct {
	svc *service.PortalService
}

func NewPortalTemplateHandler(svc *service.PortalService) *PortalTemplateHandler {
	return &PortalTemplateHandler{svc: svc}
}

func (h *PortalTemplateHandler) ListTemplates(c echo.Context) error {
	wid := c.Param("wid")
	list, err := h.svc.ListTemplates(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list templates", "error", err)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to list templates")
	}
	return c.JSON(http.StatusOK, list)
}

func (h *PortalTemplateHandler) CreateTemplate(c echo.Context) error {
	wid := c.Param("wid")
	var body dto.CreateTemplateRequest
	if err := c.Bind(&body); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to bind template create body", "error", err, "workspace_id", wid)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	t := body.ToDomain(wid)
	if err := h.svc.CreateTemplate(c.Request().Context(), t); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to create template", "error", err, "workspace_id", wid, "unique_key", body.UniqueKey)
		return safeHTTPError(err, http.StatusBadRequest, "failed to create template")
	}
	return c.JSON(http.StatusCreated, t)
}

func (h *PortalTemplateHandler) PatchTemplate(c echo.Context) error {
	wid := c.Param("wid")
	tid := c.Param("tid")
	var body dto.PatchTemplateRequest
	if err := c.Bind(&body); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to bind template patch body", "error", err, "workspace_id", wid, "template_id", tid)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if err := h.svc.PatchTemplate(c.Request().Context(), wid, tid, body.ToPatch()); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to patch template", "error", err, "workspace_id", wid, "template_id", tid)
		return safeHTTPError(err, http.StatusBadRequest, "failed to patch template")
	}
	t, err := h.svc.GetTemplate(c.Request().Context(), tid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to load template after patch", "error", err, "workspace_id", wid, "template_id", tid)
		return safeHTTPError(err, http.StatusNotFound, "template not found")
	}
	return c.JSON(http.StatusOK, t)
}

func (h *PortalTemplateHandler) DeleteTemplate(c echo.Context) error {
	wid := c.Param("wid")
	tid := c.Param("tid")
	if err := h.svc.DeleteTemplate(c.Request().Context(), wid, tid); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to delete template", "error", err, "workspace_id", wid, "template_id", tid)
		return safeHTTPError(err, http.StatusBadRequest, "failed to delete template")
	}
	return c.NoContent(http.StatusNoContent)
}
