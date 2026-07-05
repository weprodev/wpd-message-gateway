// Package service contains application services that coordinate domain logic.
// Services depend only on port interfaces — they must never import infrastructure.
package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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
	Gate         port.AuthorizationGate
	TxManager    port.TransactionManager
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
	gate         port.AuthorizationGate
	txManager    port.TransactionManager

	// jwtSecret and jwtTTL are unexported; accessed only within this package.
	jwtSecret string
	jwtTTL    time.Duration
}

// NewPortalService constructs a PortalService from its dependencies.
// If deps.JWTTTL is zero, it defaults to 24 hours.
func NewPortalService(deps PortalDeps) *PortalService {
	if deps.JWTTTL <= 0 {
		deps.JWTTTL = 24 * time.Hour
	}
	txM := deps.TxManager
	if txM == nil {
		txM = &nopTxManager{}
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
		gate:         deps.Gate,
		txManager:    txM,
		jwtSecret:    deps.JWTSecret,
		jwtTTL:       deps.JWTTTL,
	}
}

type nopTxManager struct{}

func (m *nopTxManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (s *PortalService) UserByID(ctx context.Context, id string) (*domain.User, error) {
	u, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}

// RequireMember asserts that userID is a member of workspaceID and returns their role.
// Returns port.ErrNotFound if the user is not a member.
func (s *PortalService) RequireMember(ctx context.Context, workspaceID, userID string) (string, error) {
	return s.members.GetRole(ctx, workspaceID, userID)
}

// RequireAdmin asserts that userID holds the admin role in workspaceID.
// Returns port.ErrUnauthorized (wrapped) if the role is insufficient.
func (s *PortalService) RequireAdmin(ctx context.Context, workspaceID, userID string) error {
	role, err := s.members.GetRole(ctx, workspaceID, userID)
	if err != nil {
		return err
	}
	if role != domain.RoleAdmin {
		return fmt.Errorf("admin role required: %w", port.ErrUnauthorized)
	}
	return nil
}

// CreateWorkspace creates a new workspace owned by userID, adds userID as admin, and assigns
// the admin role in the RBAC gate. Returns port.ErrInvalidInput if name or slug is empty.
func (s *PortalService) CreateWorkspace(ctx context.Context, userID, name, slug, iconKey string) (*domain.Workspace, error) {
	slog.InfoContext(ctx, "creating workspace", "user_id", userID, "name", name, "slug", slug)
	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(strings.ToLower(slug))
	if name == "" || slug == "" {
		slog.WarnContext(ctx, "workspace creation failed: invalid input", "user_id", userID)
		return nil, fmt.Errorf("name and slug required: %w", port.ErrInvalidInput)
	}
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "workspace creation failed: user not found", "error", err, "user_id", userID)
		return nil, err
	}
	var w *domain.Workspace
	err = s.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {
		w = &domain.Workspace{
			Name:       name,
			Slug:       slug,
			AdminEmail: u.Email,
			Status:     "active",
			Visibility: "private",
			IconKey:    strings.TrimSpace(iconKey),
		}
		if err := s.workspaces.Create(txCtx, w, userID); err != nil {
			return err
		}
		if err := s.members.Add(txCtx, w.ID, userID, domain.RoleAdmin); err != nil {
			return err
		}
		if err := s.gate.AssignRole(txCtx, "users", userID, w.ID, domain.RoleAdmin); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		slog.ErrorContext(ctx, "workspace creation failed", "error", err, "user_id", userID)
		return nil, err
	}
	slog.InfoContext(ctx, "workspace created successfully", "workspace_id", w.ID, "user_id", userID)
	return w, nil
}

func (s *PortalService) JoinWorkspaceWithPIN(ctx context.Context, userID, slug, pin string) error {
	slog.InfoContext(ctx, "joining workspace with PIN", "user_id", userID, "slug", slug)
	slug = strings.TrimSpace(strings.ToLower(slug))
	ws, err := s.workspaces.GetBySlug(ctx, slug)
	if err != nil || ws == nil {
		slog.WarnContext(ctx, "join workspace with PIN failed: workspace not found", "slug", slug)
		return errors.New("workspace not found")
	}
	if ws.HashedPin == "" {
		slog.WarnContext(ctx, "join workspace with PIN failed: PIN join not configured", "workspace_id", ws.ID)
		return errors.New("workspace does not use PIN join")
	}
	if !crypto.CheckSecretHash(pin, ws.HashedPin) {
		slog.WarnContext(ctx, "join workspace with PIN failed: invalid PIN", "workspace_id", ws.ID, "user_id", userID)
		return errors.New("invalid PIN")
	}
	err = s.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := s.members.Add(txCtx, ws.ID, userID, domain.RoleMember); err != nil {
			return err
		}
		if err := s.gate.AssignRole(txCtx, "users", userID, ws.ID, domain.RoleMember); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		slog.ErrorContext(ctx, "join workspace with PIN failed: database error", "error", err, "workspace_id", ws.ID, "user_id", userID)
		return err
	}
	slog.InfoContext(ctx, "joined workspace successfully with PIN", "workspace_id", ws.ID, "user_id", userID)
	return nil
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
	workspaces, err := s.workspaces.ListForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	memberWorkspaceIDs := make([]string, 0, len(workspaces))
	for i := range workspaces {
		w := &workspaces[i]
		if w.Role != "" {
			memberWorkspaceIDs = append(memberWorkspaceIDs, w.ID)
			continue
		}
		if w.Visibility == "public" {
			w.Role = domain.RoleViewer
			w.Permissions = append([]string(nil), domain.PublicGuestPermissions...)
		}
	}

	permsByTeam, err := s.gate.GetPermissionsForTeams(ctx, "users", userID, memberWorkspaceIDs)
	if err != nil {
		return nil, fmt.Errorf("load workspace permissions: %w", err)
	}

	for i := range workspaces {
		w := &workspaces[i]
		if w.Role == "" {
			continue
		}
		if perms, ok := permsByTeam[w.ID]; ok {
			w.Permissions = perms
		}
	}

	return workspaces, nil
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
	// Fetch current role to remove the exact RBAC assignment.
	role, err := s.members.GetRole(ctx, workspaceID, userID)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			// No membership record; nothing to remove.
			return nil
		}
		return fmt.Errorf("wpd-message-gateway: get member role: %w", err)
	}

	return s.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := s.gate.RemoveRole(txCtx, "users", userID, workspaceID, role); err != nil {
			return fmt.Errorf("remove role %s: %w", role, err)
		}
		if err := s.members.Remove(txCtx, workspaceID, userID); err != nil {
			return err
		}
		return nil
	})
}

func (s *PortalService) ListAPIKeys(ctx context.Context, workspaceID string) ([]domain.APIKey, error) {
	return s.apiKeys.ListByWorkspace(ctx, workspaceID)
}

func (s *PortalService) CreateAPIKey(ctx context.Context, workspaceID, name string) (*domain.APIKey, string, error) {
	slog.InfoContext(ctx, "creating api key", "workspace_id", workspaceID, "name", name)
	clientID, secret, hash, err := s.RandomAPIKeyCredentials()
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate api key credentials", "error", err, "workspace_id", workspaceID)
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
		slog.ErrorContext(ctx, "failed to persist api key in database", "error", err, "workspace_id", workspaceID)
		return nil, "", err
	}
	slog.InfoContext(ctx, "api key created successfully", "workspace_id", workspaceID, "api_key_id", k.ID)
	return k, secret, nil
}

func (s *PortalService) DeleteAPIKey(ctx context.Context, workspaceID, keyID string) error {
	slog.InfoContext(ctx, "deleting api key", "workspace_id", workspaceID, "api_key_id", keyID)
	k, err := s.apiKeys.GetByID(ctx, keyID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch api key for deletion", "error", err, "workspace_id", workspaceID, "api_key_id", keyID)
		return err
	}
	if k.WorkspaceID != workspaceID {
		slog.WarnContext(ctx, "api key delete rejected: workspace mismatch", "workspace_id", workspaceID, "api_key_id", keyID, "actual_workspace_id", k.WorkspaceID)
		return errors.New("api key not in workspace")
	}
	if err := s.apiKeys.Delete(ctx, keyID); err != nil {
		slog.ErrorContext(ctx, "failed to delete api key from database", "error", err, "workspace_id", workspaceID, "api_key_id", keyID)
		return err
	}
	slog.InfoContext(ctx, "api key deleted successfully", "workspace_id", workspaceID, "api_key_id", keyID)
	return nil
}

func (s *PortalService) RegenerateAPIKey(ctx context.Context, workspaceID, keyID string) (string, error) {
	slog.InfoContext(ctx, "regenerating api key", "workspace_id", workspaceID, "api_key_id", keyID)
	k, err := s.apiKeys.GetByID(ctx, keyID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch api key for regeneration", "error", err, "workspace_id", workspaceID, "api_key_id", keyID)
		return "", err
	}
	if k.WorkspaceID != workspaceID {
		slog.WarnContext(ctx, "api key regenerate rejected: workspace mismatch", "workspace_id", workspaceID, "api_key_id", keyID, "actual_workspace_id", k.WorkspaceID)
		return "", errors.New("api key not in workspace")
	}
	_, secret, hash, err := s.RandomAPIKeyCredentials()
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate credentials for api key regeneration", "error", err, "workspace_id", workspaceID, "api_key_id", keyID)
		return "", err
	}
	if err := s.apiKeys.UpdateSecret(ctx, keyID, k.ClientID, hash); err != nil {
		slog.ErrorContext(ctx, "failed to update api key secret in database", "error", err, "workspace_id", workspaceID, "api_key_id", keyID)
		return "", err
	}
	slog.InfoContext(ctx, "api key regenerated successfully", "workspace_id", workspaceID, "api_key_id", keyID)
	return secret, nil
}

func (s *PortalService) ListLogs(ctx context.Context, q port.MessageLogQuery) ([]domain.MessageRequestLogWithSource, int, error) {
	return s.logs.ListWithSource(ctx, q)
}

func (s *PortalService) ListIntegrations(ctx context.Context, workspaceID string) ([]domain.Integration, error) {
	list, err := s.integrations.ListByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Integration, 0, len(list))
	for _, intg := range list {
		if intg.ProviderName != "memory" {
			out = append(out, intg)
		}
	}
	return out, nil
}

func (s *PortalService) UpsertIntegration(ctx context.Context, intg *domain.Integration) error {
	slog.InfoContext(ctx, "upserting integration", "workspace_id", intg.WorkspaceID, "provider_name", intg.ProviderName, "channel_type", intg.ChannelType)
	err := s.integrations.Upsert(ctx, intg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to upsert integration", "error", err, "workspace_id", intg.WorkspaceID, "provider_name", intg.ProviderName, "channel_type", intg.ChannelType)
	} else {
		slog.InfoContext(ctx, "integration upserted successfully", "workspace_id", intg.WorkspaceID, "provider_name", intg.ProviderName, "channel_type", intg.ChannelType, "integration_id", intg.ID)
	}
	return err
}

func (s *PortalService) DeleteIntegration(ctx context.Context, workspaceID, integrationID string) error {
	slog.InfoContext(ctx, "deleting integration", "workspace_id", workspaceID, "integration_id", integrationID)
	intg, err := s.integrations.GetByID(ctx, integrationID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch integration for deletion", "error", err, "workspace_id", workspaceID, "integration_id", integrationID)
		return err
	}
	if intg.WorkspaceID != workspaceID {
		slog.WarnContext(ctx, "integration delete rejected: workspace mismatch", "workspace_id", workspaceID, "integration_id", integrationID, "actual_workspace_id", intg.WorkspaceID)
		return errors.New("integration not in workspace")
	}
	if err := s.integrations.Delete(ctx, integrationID); err != nil {
		slog.ErrorContext(ctx, "failed to delete integration from database", "error", err, "workspace_id", workspaceID, "integration_id", integrationID)
		return err
	}
	slog.InfoContext(ctx, "integration deleted successfully", "workspace_id", workspaceID, "integration_id", integrationID)
	return nil
}

func (s *PortalService) ListTemplates(ctx context.Context, workspaceID string) ([]domain.Template, error) {
	return s.templates.ListByWorkspace(ctx, workspaceID)
}

func (s *PortalService) CreateTemplate(ctx context.Context, t *domain.Template) error {
	slog.InfoContext(ctx, "creating template", "workspace_id", t.WorkspaceID, "unique_key", t.UniqueKey, "channel_type", t.ChannelType)
	err := s.templates.Create(ctx, t)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create template", "error", err, "workspace_id", t.WorkspaceID, "unique_key", t.UniqueKey)
	} else {
		slog.InfoContext(ctx, "template created successfully", "workspace_id", t.WorkspaceID, "template_id", t.ID, "unique_key", t.UniqueKey)
	}
	return err
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
	slog.InfoContext(ctx, "patching template", "workspace_id", workspaceID, "template_id", templateID)
	t, err := s.templates.GetByID(ctx, templateID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch template for patching", "error", err, "workspace_id", workspaceID, "template_id", templateID)
		return err
	}
	if t.WorkspaceID != workspaceID {
		slog.WarnContext(ctx, "template patch rejected: workspace mismatch", "workspace_id", workspaceID, "template_id", templateID, "actual_workspace_id", t.WorkspaceID)
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
	if err := s.templates.Update(ctx, t); err != nil {
		slog.ErrorContext(ctx, "failed to update patched template in database", "error", err, "workspace_id", workspaceID, "template_id", templateID)
		return err
	}
	slog.InfoContext(ctx, "template patched successfully", "workspace_id", workspaceID, "template_id", templateID)
	return nil
}

func (s *PortalService) GetTemplate(ctx context.Context, templateID string) (*domain.Template, error) {
	return s.templates.GetByID(ctx, templateID)
}

func (s *PortalService) DeleteTemplate(ctx context.Context, workspaceID, templateID string) error {
	slog.InfoContext(ctx, "deleting template", "workspace_id", workspaceID, "template_id", templateID)
	t, err := s.templates.GetByID(ctx, templateID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch template for deletion", "error", err, "workspace_id", workspaceID, "template_id", templateID)
		return err
	}
	if t.WorkspaceID != workspaceID {
		slog.WarnContext(ctx, "template delete rejected: workspace mismatch", "workspace_id", workspaceID, "template_id", templateID, "actual_workspace_id", t.WorkspaceID)
		return errors.New("template not in workspace")
	}
	if err := s.templates.Delete(ctx, templateID); err != nil {
		slog.ErrorContext(ctx, "failed to delete template from database", "error", err, "workspace_id", workspaceID, "template_id", templateID)
		return err
	}
	slog.InfoContext(ctx, "template deleted successfully", "workspace_id", workspaceID, "template_id", templateID)
	return nil
}

func (s *PortalService) GetSettings(ctx context.Context, workspaceID string) (map[string]string, error) {
	return s.settings.GetAll(ctx, workspaceID)
}

func (s *PortalService) PatchSettings(ctx context.Context, workspaceID string, kv map[string]string) error {
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	slog.InfoContext(ctx, "patching workspace settings", "workspace_id", workspaceID, "keys", keys)

	for k, v := range kv {
		if err := domain.ValidateWorkspaceSettingValue(k, v); err != nil {
			slog.WarnContext(ctx, "settings patch rejected: invalid value",
				"workspace_id", workspaceID,
				"key", k,
				"error", err,
			)
			return fmt.Errorf("%w: %w", port.ErrInvalidInput, err)
		}

		if k == domain.SettingKeyMessageDispatchMode && v == string(domain.DispatchProvider) {
			list, err := s.integrations.ListByWorkspace(ctx, workspaceID)
			if err != nil {
				return fmt.Errorf("failed to check integrations for provider enablement: %w", err)
			}
			hasConnected := false
			for _, intg := range list {
				if intg.ProviderName != domain.ProviderNameMemory && intg.Status == domain.IntegrationStatusConnected {
					hasConnected = true
					break
				}
			}
			if !hasConnected {
				return fmt.Errorf("%w: configure a provider in Integrations before switching to provider mode", port.ErrInvalidInput)
			}
		}

		if err := s.settings.Set(ctx, workspaceID, k, v); err != nil {
			return err
		}
	}

	slog.InfoContext(ctx, "workspace settings patched successfully", "workspace_id", workspaceID, "keys", keys)
	return nil
}

const invitationTTL = 7 * 24 * time.Hour

// NewPendingInvitation builds a workspace invitation with the standard expiry window.
func (s *PortalService) NewPendingInvitation(workspaceID, email, role string) *domain.Invitation {
	return &domain.Invitation{
		WorkspaceID: workspaceID,
		Email:       strings.ToLower(strings.TrimSpace(email)),
		Role:        role,
		ExpiresAt:   time.Now().Add(invitationTTL),
		Status:      "pending",
	}
}

// ListInvitations returns all pending invitations for a workspace.
func (s *PortalService) ListInvitations(ctx context.Context, workspaceID string) ([]domain.Invitation, error) {
	return s.invites.ListPendingByWorkspace(ctx, workspaceID)
}

// CreateInvitation generates a secure invitation token, stores its hash, and persists the
// invitation. The returned plaintext token must be sent to the invitee and is never stored.
func (s *PortalService) CreateInvitation(ctx context.Context, inv *domain.Invitation) (rawToken string, err error) {
	if !domain.IsWorkspaceRole(inv.Role) {
		return "", fmt.Errorf("invalid role %q: %w", inv.Role, port.ErrInvalidInput)
	}

	inv.Email = strings.ToLower(strings.TrimSpace(inv.Email))
	if inv.Email == "" {
		return "", fmt.Errorf("email required: %w", port.ErrInvalidInput)
	}

	isMember, err := s.members.MemberExistsByEmail(ctx, inv.WorkspaceID, inv.Email)
	if err != nil {
		return "", err
	}
	if isMember {
		return "", fmt.Errorf("user is already a member of this workspace: %w", port.ErrInvalidInput)
	}

	hasPending, err := s.invites.PendingInvitationExistsByEmail(ctx, inv.WorkspaceID, inv.Email)
	if err != nil {
		return "", err
	}
	if hasPending {
		return "", fmt.Errorf("a pending invitation already exists for this email: %w", port.ErrInvalidInput)
	}

	// Generate a cryptographically random token from two UUIDs.
	rawToken = uuid.NewString() + uuid.NewString()
	sum := sha256.Sum256([]byte(rawToken))
	inv.TokenHash = hex.EncodeToString(sum[:])
	return rawToken, s.invites.Create(ctx, inv)
}

// JWTSecret exposes the portal JWT secret for auth-layer consumers within this package.
func (s *PortalService) JWTSecret() string { return s.jwtSecret }

// JWTTTL exposes the portal JWT TTL for auth-layer consumers within this package.
func (s *PortalService) JWTTTL() time.Duration { return s.jwtTTL }

// GetProviderFields returns the configuration fields schema for the named provider.
func (s *PortalService) GetProviderFields(ctx context.Context, providerName string) ([]domain.ProviderConfigField, error) {
	slog.InfoContext(ctx, "fetching provider fields schema", "provider_name", providerName)
	if providerName == "memory" {
		return []domain.ProviderConfigField{}, nil
	}
	fields, err := s.integrations.GetProviderFields(ctx, providerName)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch provider fields schema", "error", err, "provider_name", providerName)
		return nil, err
	}
	return fields, nil
}

// ListProviders returns all available integration providers in the catalog (excluding memory).
func (s *PortalService) ListProviders(ctx context.Context) ([]domain.Provider, error) {
	slog.InfoContext(ctx, "listing integration providers")
	list, err := s.integrations.ListProviders(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list integration providers", "error", err)
		return nil, err
	}
	out := make([]domain.Provider, 0, len(list))
	for _, p := range list {
		if p.Name != "memory" {
			out = append(out, p)
		}
	}
	return out, nil
}
