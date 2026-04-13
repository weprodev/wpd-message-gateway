// Package service contains application services that coordinate domain logic.
// Services depend only on port interfaces — they must never import infrastructure.
package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/go-pkg/crypto"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

// PortalDeps groups all dependencies for PortalService.
// Using a struct prevents mis-ordering of arguments and makes it obvious
// which dependency is which at the call site.
type PortalDeps struct {
	Users        port.UserRepository
	Members      port.WorkspaceMemberRepository
	Workspaces   port.WorkspaceRepository
	APIKeys      port.APIKeyRepository
	Integrations port.IntegrationRepository
	Templates    port.TemplateRepository
	Logs         port.MessageRequestLogRepository
	Invitations  port.InvitationRepository
	Settings     port.WorkspaceSettingsRepository
	JWTSecret    string
	JWTTTL       time.Duration
}

// PortalService handles portal authentication and workspace-scoped operations.
type PortalService struct {
	users        port.UserRepository
	members      port.WorkspaceMemberRepository
	workspaces   port.WorkspaceRepository
	apiKeys      port.APIKeyRepository
	integrations port.IntegrationRepository
	templates    port.TemplateRepository
	logs         port.MessageRequestLogRepository
	invites      port.InvitationRepository
	settings     port.WorkspaceSettingsRepository

	JWTSecret string
	JWTTTL    time.Duration
}

// NewPortalService constructs a PortalService from its dependencies.
// If deps.JWTTTL is zero, it defaults to 24 hours.
func NewPortalService(deps PortalDeps) *PortalService {
	if deps.JWTTTL <= 0 {
		deps.JWTTTL = 24 * time.Hour
	}
	return &PortalService{
		users:        deps.Users,
		members:      deps.Members,
		workspaces:   deps.Workspaces,
		apiKeys:      deps.APIKeys,
		integrations: deps.Integrations,
		templates:    deps.Templates,
		logs:         deps.Logs,
		invites:      deps.Invitations,
		settings:     deps.Settings,
		JWTSecret:    deps.JWTSecret,
		JWTTTL:       deps.JWTTTL,
	}
}

func (s *PortalService) UserByID(ctx context.Context, id string) (*domain.User, error) {
	u, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}

func (s *PortalService) RequireMember(ctx context.Context, workspaceID, userID string) (string, error) {
	return s.members.GetRole(ctx, workspaceID, userID)
}

func (s *PortalService) RequireAdmin(ctx context.Context, workspaceID, userID string) error {
	role, err := s.members.GetRole(ctx, workspaceID, userID)
	if err != nil {
		return err
	}
	if role != "admin" {
		return errors.New("admin role required")
	}
	return nil
}

func (s *PortalService) CreateWorkspace(ctx context.Context, userID, name, uniqueKey, iconKey string) (*domain.Workspace, error) {
	name = strings.TrimSpace(name)
	uniqueKey = strings.TrimSpace(strings.ToLower(uniqueKey))
	if name == "" || uniqueKey == "" {
		return nil, errors.New("name and unique_key required")
	}
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	w := &domain.Workspace{
		Name:       name,
		UniqueKey:  uniqueKey,
		AdminEmail: u.Email,
		Status:     "active",
		Visibility: "private",
		IconKey:    strings.TrimSpace(iconKey),
	}
	if err := s.workspaces.Create(ctx, w); err != nil {
		return nil, err
	}
	if err := s.members.Add(ctx, w.ID, userID, "admin"); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *PortalService) JoinWorkspaceWithPIN(ctx context.Context, userID, uniqueKey, pin string) error {
	uniqueKey = strings.TrimSpace(strings.ToLower(uniqueKey))
	ws, err := s.workspaces.GetByUniqueKey(ctx, uniqueKey)
	if err != nil || ws == nil {
		return errors.New("workspace not found")
	}
	if ws.HashedPin == "" {
		return errors.New("workspace does not use PIN join")
	}
	if !crypto.CheckSecretHash(pin, ws.HashedPin) {
		return errors.New("invalid PIN")
	}
	return s.members.Add(ctx, ws.ID, userID, "member")
}

func (s *PortalService) RandomAPIKeyCredentials() (clientID, secret string, secretHash string, err error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", "", err
	}
	secret = base64.RawURLEncoding.EncodeToString(raw)
	clientID = "wk_" + strings.ReplaceAll(uuid.New().String(), "-", "")
	hash, err := crypto.HashSecret(secret)
	if err != nil {
		return "", "", "", err
	}
	return clientID, secret, hash, nil
}

// IntegrationConfigJSON parses or builds integration config from JSON bytes for API handlers.
func IntegrationConfigJSON(b []byte) ([]byte, error) {
	if len(b) == 0 {
		return []byte("{}"), nil
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

func (s *PortalService) ListWorkspaces(ctx context.Context, userID string) ([]domain.Workspace, error) {
	return s.workspaces.ListForUser(ctx, userID)
}

func (s *PortalService) WorkspaceByID(ctx context.Context, id string) (*domain.Workspace, error) {
	return s.workspaces.GetByID(ctx, id)
}

// WorkspacePatch carries optional updates for PATCH workspace.
type WorkspacePatch struct {
	Name       *string
	Visibility *string
	IconKey    *string
}

func (s *PortalService) PatchWorkspace(ctx context.Context, id string, p WorkspacePatch) error {
	w, err := s.workspaces.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p.Name != nil {
		w.Name = strings.TrimSpace(*p.Name)
	}
	if p.Visibility != nil {
		w.Visibility = *p.Visibility
	}
	if p.IconKey != nil {
		w.IconKey = strings.TrimSpace(*p.IconKey)
	}
	return s.workspaces.Update(ctx, w)
}

func (s *PortalService) ListMembers(ctx context.Context, workspaceID string) ([]domain.WorkspaceMember, error) {
	return s.members.ListMembers(ctx, workspaceID)
}

func (s *PortalService) RemoveMember(ctx context.Context, workspaceID, userID string) error {
	return s.members.Remove(ctx, workspaceID, userID)
}

func (s *PortalService) ListAPIKeys(ctx context.Context, workspaceID string) ([]domain.APIKey, error) {
	return s.apiKeys.ListByWorkspace(ctx, workspaceID)
}

func (s *PortalService) CreateAPIKey(ctx context.Context, workspaceID, name string) (*domain.APIKey, string, error) {
	clientID, secret, hash, err := s.RandomAPIKeyCredentials()
	if err != nil {
		return nil, "", err
	}
	k := &domain.APIKey{
		WorkspaceID:      workspaceID,
		ClientID:         clientID,
		ClientSecretHash: hash,
		Name:             strings.TrimSpace(name),
		IsActive:         true,
	}
	if err := s.apiKeys.Create(ctx, k); err != nil {
		return nil, "", err
	}
	return k, secret, nil
}

func (s *PortalService) DeleteAPIKey(ctx context.Context, workspaceID, keyID string) error {
	k, err := s.apiKeys.GetByID(ctx, keyID)
	if err != nil {
		return err
	}
	if k.WorkspaceID != workspaceID {
		return errors.New("api key not in workspace")
	}
	return s.apiKeys.Delete(ctx, keyID)
}

func (s *PortalService) RegenerateAPIKey(ctx context.Context, workspaceID, keyID string) (string, error) {
	k, err := s.apiKeys.GetByID(ctx, keyID)
	if err != nil {
		return "", err
	}
	if k.WorkspaceID != workspaceID {
		return "", errors.New("api key not in workspace")
	}
	_, secret, hash, err := s.RandomAPIKeyCredentials()
	if err != nil {
		return "", err
	}
	if err := s.apiKeys.UpdateSecret(ctx, keyID, k.ClientID, hash); err != nil {
		return "", err
	}
	return secret, nil
}

func (s *PortalService) ListLogs(ctx context.Context, q port.MessageLogQuery) ([]domain.MessageRequestLogWithSource, int, error) {
	return s.logs.ListWithSource(ctx, q)
}

func (s *PortalService) ListIntegrations(ctx context.Context, workspaceID string) ([]domain.Integration, error) {
	return s.integrations.ListByWorkspace(ctx, workspaceID)
}

func (s *PortalService) UpsertIntegration(ctx context.Context, intg *domain.Integration) error {
	return s.integrations.Upsert(ctx, intg)
}

func (s *PortalService) DeleteIntegration(ctx context.Context, workspaceID, integrationID string) error {
	intg, err := s.integrations.GetByID(ctx, integrationID)
	if err != nil {
		return err
	}
	if intg.WorkspaceID != workspaceID {
		return errors.New("integration not in workspace")
	}
	return s.integrations.Delete(ctx, integrationID)
}

func (s *PortalService) ListTemplates(ctx context.Context, workspaceID string) ([]domain.Template, error) {
	return s.templates.ListByWorkspace(ctx, workspaceID)
}

func (s *PortalService) CreateTemplate(ctx context.Context, t *domain.Template) error {
	return s.templates.Create(ctx, t)
}

// TemplatePatch carries optional fields for PATCH template.
type TemplatePatch struct {
	Name        *string
	UniqueKey   *string
	ChannelType *string
	Category    *string
	Subject     *string
	ContentHTML *string
	IsActive    *bool
	IsDefault   *bool
}

func (s *PortalService) PatchTemplate(ctx context.Context, workspaceID, templateID string, p TemplatePatch) error {
	t, err := s.templates.GetByID(ctx, templateID)
	if err != nil {
		return err
	}
	if t.WorkspaceID != workspaceID {
		return errors.New("template not in workspace")
	}
	if p.Name != nil {
		t.Name = strings.TrimSpace(*p.Name)
	}
	if p.UniqueKey != nil {
		t.UniqueKey = strings.TrimSpace(*p.UniqueKey)
	}
	if p.ChannelType != nil {
		t.ChannelType = *p.ChannelType
	}
	if p.Category != nil {
		t.Category = *p.Category
	}
	if p.Subject != nil {
		t.Subject = *p.Subject
	}
	if p.ContentHTML != nil {
		t.ContentHTML = *p.ContentHTML
	}
	if p.IsActive != nil {
		t.IsActive = *p.IsActive
	}
	if p.IsDefault != nil {
		t.IsDefault = *p.IsDefault
	}
	return s.templates.Update(ctx, t)
}

func (s *PortalService) GetTemplate(ctx context.Context, templateID string) (*domain.Template, error) {
	return s.templates.GetByID(ctx, templateID)
}

func (s *PortalService) DeleteTemplate(ctx context.Context, workspaceID, templateID string) error {
	t, err := s.templates.GetByID(ctx, templateID)
	if err != nil {
		return err
	}
	if t.WorkspaceID != workspaceID {
		return errors.New("template not in workspace")
	}
	return s.templates.Delete(ctx, templateID)
}

func (s *PortalService) GetSettings(ctx context.Context, workspaceID string) (map[string]string, error) {
	return s.settings.GetAll(ctx, workspaceID)
}

func (s *PortalService) PatchSettings(ctx context.Context, workspaceID string, kv map[string]string) error {
	for k, v := range kv {
		if err := s.settings.Set(ctx, workspaceID, k, v); err != nil {
			return err
		}
	}
	return nil
}

func (s *PortalService) ListInvitations(ctx context.Context, workspaceID string) ([]domain.Invitation, error) {
	return s.invites.ListPendingByWorkspace(ctx, workspaceID)
}

func (s *PortalService) CreateInvitation(ctx context.Context, inv *domain.Invitation) error {
	return s.invites.Create(ctx, inv)
}
