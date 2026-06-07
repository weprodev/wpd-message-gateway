package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	customMiddleware "github.com/weprodev/wpd-message-gateway/internal/presentation/middleware"
)

type PortalWorkspaceHandler struct {
	svc *service.PortalService
}

func NewPortalWorkspaceHandler(svc *service.PortalService) *PortalWorkspaceHandler {
	return &PortalWorkspaceHandler{svc: svc}
}

type createWorkspaceBody struct {
	Name      string `json:"name"`
	UniqueKey string `json:"unique_key"`
	IconKey   string `json:"icon_key"`
}

type joinBody struct {
	UniqueKey string `json:"unique_key"`
	PIN       string `json:"pin"`
}

type patchWorkspaceBody struct {
	Name       *string `json:"name"`
	Visibility *string `json:"visibility"`
	IconKey    *string `json:"icon_key"`
}

type invitationBody struct {
	Email string `json:"email"`
	Role  string `json:"role"`
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
	var body createWorkspaceBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	w, err := h.svc.CreateWorkspace(c.Request().Context(), uid, body.Name, body.UniqueKey, body.IconKey)
	if err != nil {
		return safeHTTPError(err, http.StatusBadRequest, "failed to create workspace")
	}
	return c.JSON(http.StatusCreated, w)
}

func (h *PortalWorkspaceHandler) JoinWorkspace(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	var body joinBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if err := h.svc.JoinWorkspaceWithPIN(c.Request().Context(), uid, body.UniqueKey, body.PIN); err != nil {
		return safeHTTPError(err, http.StatusBadRequest, "failed to join workspace")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalWorkspaceHandler) GetWorkspace(c echo.Context) error {
	wid := c.Param("wid")
	w, err := h.svc.WorkspaceByID(c.Request().Context(), wid)
	if err != nil {
		return safeHTTPError(err, http.StatusNotFound, "workspace not found")
	}
	return c.JSON(http.StatusOK, w)
}

func (h *PortalWorkspaceHandler) PatchWorkspace(c echo.Context) error {
	wid := c.Param("wid")
	var body patchWorkspaceBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if err := h.svc.PatchWorkspace(c.Request().Context(), wid, service.WorkspacePatch{
		Name: body.Name, Visibility: body.Visibility, IconKey: body.IconKey,
	}); err != nil {
		return safeHTTPError(err, http.StatusBadRequest, "failed to patch workspace")
	}
	w, err := h.svc.WorkspaceByID(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to load workspace after patch", "error", err)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to load workspace")
	}
	return c.JSON(http.StatusOK, w)
}

func (h *PortalWorkspaceHandler) ListMembers(c echo.Context) error {
	wid := c.Param("wid")
	members, err := h.svc.ListMembers(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list members", "error", err)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to list members")
	}
	return c.JSON(http.StatusOK, members)
}

func (h *PortalWorkspaceHandler) RemoveMember(c echo.Context) error {
	wid := c.Param("wid")
	target := c.Param("userId")
	if err := h.svc.RemoveMember(c.Request().Context(), wid, target); err != nil {
		return safeHTTPError(err, http.StatusBadRequest, "failed to remove member")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalWorkspaceHandler) ListInvitations(c echo.Context) error {
	wid := c.Param("wid")
	list, err := h.svc.ListInvitations(c.Request().Context(), wid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list invitations", "error", err)
		return safeHTTPError(err, http.StatusInternalServerError, "failed to list invitations")
	}
	return c.JSON(http.StatusOK, list)
}

func (h *PortalWorkspaceHandler) CreateInvitation(c echo.Context) error {
	wid := c.Param("wid")
	var body invitationBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if body.Email == "" || body.Role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email and role required")
	}
	inv := &domain.Invitation{
		WorkspaceID: wid,
		Email:       body.Email,
		Role:        body.Role,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
		Status:      "pending",
	}
	rawToken, err := h.svc.CreateInvitation(c.Request().Context(), inv)
	if err != nil {
		return safeHTTPError(err, http.StatusBadRequest, "failed to create invitation")
	}
	return c.JSON(http.StatusCreated, map[string]any{
		"id":         inv.ID,
		"email":      inv.Email,
		"role":       inv.Role,
		"expires_at": inv.ExpiresAt,
		"token":      rawToken,
	})
}
