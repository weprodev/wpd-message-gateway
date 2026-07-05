package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/internal/presentation/dto"
	customMiddleware "github.com/weprodev/wpd-message-gateway/internal/presentation/middleware"
)

type PortalWorkspaceHandler struct {
	svc *service.PortalService
}

func NewPortalWorkspaceHandler(svc *service.PortalService) *PortalWorkspaceHandler {
	return &PortalWorkspaceHandler{svc: svc}
}

func (h *PortalWorkspaceHandler) ListWorkspaces(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	list, err := h.svc.ListWorkspaces(c.Request().Context(), uid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list workspaces", "error", err)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to list workspaces")
	}
	return c.JSON(http.StatusOK, list)
}

func (h *PortalWorkspaceHandler) CreateWorkspace(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	var body dto.CreateWorkspaceRequest
	if err := c.Bind(&body); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to bind workspace create body", "error", err, "user_id", uid)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	w, err := h.svc.CreateWorkspace(c.Request().Context(), uid, body.Name, body.Slug, body.IconKey)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to create workspace", "error", err, "user_id", uid, "name", body.Name)
		if errors.Is(err, port.ErrConflict) {
			return echo.NewHTTPError(http.StatusConflict, "workspace slug already exists")
		}
		return safeHTTPError(err, http.StatusBadRequest, "failed to create workspace")
	}
	return c.JSON(http.StatusCreated, w)
}

func (h *PortalWorkspaceHandler) JoinWorkspace(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	var body dto.JoinWorkspaceRequest
	if err := c.Bind(&body); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to bind workspace join body", "error", err, "user_id", uid)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if err := h.svc.JoinWorkspaceWithPIN(c.Request().Context(), uid, body.Slug, body.PIN); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to join workspace", "error", err, "user_id", uid, "slug", body.Slug)
		return safeHTTPError(err, http.StatusBadRequest, "failed to join workspace")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalWorkspaceHandler) GetWorkspace(c echo.Context) error {
	wid := c.Param("wid")
	w, err := h.svc.WorkspaceByID(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to get workspace", "error", err, "workspace_id", wid)
		return safeHTTPError(err, http.StatusNotFound, "workspace not found")
	}
	return c.JSON(http.StatusOK, w)
}

func (h *PortalWorkspaceHandler) PatchWorkspace(c echo.Context) error {
	wid := c.Param("wid")
	var body dto.PatchWorkspaceRequest
	if err := c.Bind(&body); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to bind workspace patch body", "error", err, "workspace_id", wid)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if err := h.svc.PatchWorkspace(c.Request().Context(), wid, body.ToPatch()); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to patch workspace", "error", err, "workspace_id", wid)
		return safeHTTPError(err, http.StatusBadRequest, "failed to patch workspace")
	}
	w, err := h.svc.WorkspaceByID(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to load workspace after patch", "error", err, "workspace_id", wid)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to load workspace")
	}
	return c.JSON(http.StatusOK, w)
}

func (h *PortalWorkspaceHandler) ListMembers(c echo.Context) error {
	wid := c.Param("wid")
	members, err := h.svc.ListMembers(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list members", "error", err, "workspace_id", wid)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to list members")
	}
	if members == nil {
		members = []domain.WorkspaceMember{}
	}
	return c.JSON(http.StatusOK, members)
}

func (h *PortalWorkspaceHandler) RemoveMember(c echo.Context) error {
	wid := c.Param("wid")
	target := c.Param("userId")
	if err := h.svc.RemoveMember(c.Request().Context(), wid, target); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to remove member", "error", err, "workspace_id", wid, "target_user_id", target)
		return safeHTTPError(err, http.StatusBadRequest, "failed to remove member")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalWorkspaceHandler) ListInvitations(c echo.Context) error {
	wid := c.Param("wid")
	list, err := h.svc.ListInvitations(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list invitations", "error", err, "workspace_id", wid)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to list invitations")
	}
	if list == nil {
		list = []domain.Invitation{}
	}
	return c.JSON(http.StatusOK, list)
}

func (h *PortalWorkspaceHandler) CreateInvitation(c echo.Context) error {
	wid := c.Param("wid")
	var body dto.CreateInvitationRequest
	if err := c.Bind(&body); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to bind invitation body", "error", err, "workspace_id", wid)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if body.Email == "" || body.Role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email and role required")
	}
	inv := h.svc.NewPendingInvitation(wid, body.Email, body.Role)
	rawToken, err := h.svc.CreateInvitation(c.Request().Context(), inv)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to create invitation", "error", err, "workspace_id", wid, "email", body.Email)
		return safeHTTPError(err, http.StatusBadRequest, "failed to create invitation")
	}
	return c.JSON(http.StatusCreated, dto.CreateInvitationResponseFromDomain(*inv, rawToken))
}
