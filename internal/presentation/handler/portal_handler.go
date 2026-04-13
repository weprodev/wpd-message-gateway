package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/weprodev/go-pkg/sanitizer"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	customMiddleware "github.com/weprodev/wpd-message-gateway/internal/presentation/middleware"
)

// PortalHandler exposes REST endpoints for the web portal (JWT, not API keys).
type PortalHandler struct {
	svc    *service.PortalService
	authSvc *service.AuthService
}

func NewPortalHandler(svc *service.PortalService, authSvc *service.AuthService) *PortalHandler {
	return &PortalHandler{svc: svc, authSvc: authSvc}
}

type registerBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

func (h *PortalHandler) Register(c echo.Context) error {
	var body registerBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	u, err := h.authSvc.RegisterUser(c.Request().Context(), body.FirstName, body.LastName, body.Email, body.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, u)
}

func (h *PortalHandler) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing token")
	}

	if err := h.authSvc.VerifyEmail(c.Request().Context(), token); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *PortalHandler) Login(c echo.Context) error {
	var body loginBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}

	u, token, err := h.authSvc.Login(c.Request().Context(), body.Email, body.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, tokenResponse{Token: token, User: u})
}

func (h *PortalHandler) Me(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	u, err := h.svc.UserByID(c.Request().Context(), uid)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load user")
	}
	return c.JSON(http.StatusOK, u)
}

type createWorkspaceBody struct {
	Name      string `json:"name"`
	UniqueKey string `json:"unique_key"`
	IconKey   string `json:"icon_key"`
}

func (h *PortalHandler) ListWorkspaces(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	list, err := h.svc.ListWorkspaces(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

func (h *PortalHandler) CreateWorkspace(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	var body createWorkspaceBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	w, err := h.svc.CreateWorkspace(c.Request().Context(), uid, body.Name, body.UniqueKey, body.IconKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, w)
}

type joinBody struct {
	UniqueKey string `json:"unique_key"`
	PIN       string `json:"pin"`
}

func (h *PortalHandler) JoinWorkspace(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	var body joinBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if err := h.svc.JoinWorkspaceWithPIN(c.Request().Context(), uid, body.UniqueKey, body.PIN); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalHandler) GetWorkspace(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if _, err := h.svc.RequireMember(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "not a member")
	}
	w, err := h.svc.WorkspaceByID(c.Request().Context(), wid)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, w)
}

type patchWorkspaceBody struct {
	Name       *string `json:"name"`
	Visibility *string `json:"visibility"`
	IconKey    *string `json:"icon_key"`
}

func (h *PortalHandler) PatchWorkspace(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	var body patchWorkspaceBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if err := h.svc.PatchWorkspace(c.Request().Context(), wid, service.WorkspacePatch{
		Name: body.Name, Visibility: body.Visibility, IconKey: body.IconKey,
	}); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	w, err := h.svc.WorkspaceByID(c.Request().Context(), wid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, w)
}

func (h *PortalHandler) ListMembers(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if _, err := h.svc.RequireMember(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "not a member")
	}
	members, err := h.svc.ListMembers(c.Request().Context(), wid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, members)
}

func (h *PortalHandler) RemoveMember(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	target := c.Param("userId")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	if err := h.svc.RemoveMember(c.Request().Context(), wid, target); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
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
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if _, err := h.svc.RequireMember(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "not a member")
	}
	keys, err := h.svc.ListAPIKeys(c.Request().Context(), wid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	out := make([]apiKeyPublic, 0, len(keys))
	for _, k := range keys {
		out = append(out, toAPIKeyPublic(k))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *PortalHandler) CreateAPIKey(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	var body createAPIKeyBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if body.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name required")
	}
	k, secret, err := h.svc.CreateAPIKey(c.Request().Context(), wid, body.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	keyID := c.Param("keyId")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	if err := h.svc.DeleteAPIKey(c.Request().Context(), wid, keyID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalHandler) RegenerateAPIKey(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	keyID := c.Param("keyId")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	secret, err := h.svc.RegenerateAPIKey(c.Request().Context(), wid, keyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"client_secret": secret})
}

func (h *PortalHandler) ListLogs(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if _, err := h.svc.RequireMember(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "not a member")
	}
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]any{"items": rows, "total": total})
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

func (h *PortalHandler) ListIntegrations(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if _, err := h.svc.RequireMember(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "not a member")
	}
	list, err := h.svc.ListIntegrations(c.Request().Context(), wid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	out := make([]map[string]any, 0, len(list))
	for _, intg := range list {
		out = append(out, integrationToJSON(intg))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *PortalHandler) UpsertIntegration(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, integrationToJSON(*intg))
}

func (h *PortalHandler) DeleteIntegration(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	iid := c.Param("iid")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	if err := h.svc.DeleteIntegration(c.Request().Context(), wid, iid); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

type templateBody struct {
	Name        string `json:"name"`
	UniqueKey   string `json:"unique_key"`
	ChannelType string `json:"channel_type"`
	Category    string `json:"category"`
	Subject     string `json:"subject"`
	ContentHTML string `json:"content_html"`
	IsActive    *bool  `json:"is_active"`
	IsDefault   *bool  `json:"is_default"`
}

type patchTemplateBody struct {
	Name        *string `json:"name"`
	UniqueKey   *string `json:"unique_key"`
	ChannelType *string `json:"channel_type"`
	Category    *string `json:"category"`
	Subject     *string `json:"subject"`
	ContentHTML *string `json:"content_html"`
	IsActive    *bool   `json:"is_active"`
	IsDefault   *bool   `json:"is_default"`
}

func (h *PortalHandler) ListTemplates(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if _, err := h.svc.RequireMember(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "not a member")
	}
	list, err := h.svc.ListTemplates(c.Request().Context(), wid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

func (h *PortalHandler) CreateTemplate(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	var body templateBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	t := &domain.Template{
		WorkspaceID: wid,
		Name:        body.Name,
		UniqueKey:   body.UniqueKey,
		ChannelType: body.ChannelType,
		Category:    body.Category,
		Subject:     body.Subject,
		ContentHTML: sanitizer.SanitizeHTML(body.ContentHTML),
		IsActive:    true,
		IsDefault:   false,
	}
	if body.IsActive != nil {
		t.IsActive = *body.IsActive
	}
	if body.IsDefault != nil {
		t.IsDefault = *body.IsDefault
	}
	if t.ChannelType == "" {
		t.ChannelType = "email"
	}
	if err := h.svc.CreateTemplate(c.Request().Context(), t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, t)
}

func (h *PortalHandler) PatchTemplate(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	tid := c.Param("tid")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	var body patchTemplateBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	patch := service.TemplatePatch{
		Name: body.Name, UniqueKey: body.UniqueKey, ChannelType: body.ChannelType,
		Category: body.Category, Subject: body.Subject, ContentHTML: body.ContentHTML,
		IsActive: body.IsActive, IsDefault: body.IsDefault,
	}
	if body.ContentHTML != nil {
		clean := sanitizer.SanitizeHTML(*body.ContentHTML)
		patch.ContentHTML = &clean
	}
	if err := h.svc.PatchTemplate(c.Request().Context(), wid, tid, patch); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	t, err := h.svc.GetTemplate(c.Request().Context(), tid)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, t)
}

func (h *PortalHandler) DeleteTemplate(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	tid := c.Param("tid")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	if err := h.svc.DeleteTemplate(c.Request().Context(), wid, tid); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalHandler) GetSettings(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if _, err := h.svc.RequireMember(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "not a member")
	}
	m, err := h.svc.GetSettings(c.Request().Context(), wid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, m)
}

func (h *PortalHandler) PatchSettings(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	var body map[string]string
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if err := h.svc.PatchSettings(c.Request().Context(), wid, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

type invitationBody struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *PortalHandler) ListInvitations(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	list, err := h.svc.ListInvitations(c.Request().Context(), wid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

func (h *PortalHandler) CreateInvitation(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	wid := c.Param("wid")
	if err := h.svc.RequireAdmin(c.Request().Context(), wid, uid); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	var body invitationBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	if body.Email == "" || body.Role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email and role required")
	}
	rawToken := uuid.NewString() + uuid.NewString()
	sum := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(sum[:])
	inv := &domain.Invitation{
		WorkspaceID: wid,
		Email:       body.Email,
		Role:        body.Role,
		TokenHash:   tokenHash,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
		Status:      "pending",
	}
	if err := h.svc.CreateInvitation(c.Request().Context(), inv); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, map[string]any{
		"id":         inv.ID,
		"email":      inv.Email,
		"role":       inv.Role,
		"expires_at": inv.ExpiresAt,
		"token":      rawToken,
	})
}
