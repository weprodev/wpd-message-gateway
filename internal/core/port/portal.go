package port

import (
	"context"
	"time"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

// WorkspaceMemberRepository manages workspace membership.
type WorkspaceMemberRepository interface {
	Add(ctx context.Context, workspaceID, userID, role string) error
	Remove(ctx context.Context, workspaceID, userID string) error
	GetRole(ctx context.Context, workspaceID, userID string) (string, error)
	ListMembers(ctx context.Context, workspaceID string) ([]domain.WorkspaceMember, error)
}

// InvitationRepository stores pending invites.
type InvitationRepository interface {
	Create(ctx context.Context, inv *domain.Invitation) error
	ListPendingByWorkspace(ctx context.Context, workspaceID string) ([]domain.Invitation, error)
}

// MessageLogQuery filters for listing request logs.
type MessageLogQuery struct {
	WorkspaceID string
	ChannelType string
	Limit       int
	Offset      int
	From        *time.Time
	To          *time.Time
}

// WorkspaceSettingsRepository is key-value settings per workspace.
type WorkspaceSettingsRepository interface {
	Get(ctx context.Context, workspaceID, key string) (string, error)
	Set(ctx context.Context, workspaceID, key, value string) error
	GetAll(ctx context.Context, workspaceID string) (map[string]string, error)
}

// AuthorizationGate abstracts role and permission assignments.
type AuthorizationGate interface {
	AssignRole(ctx context.Context, modelType, modelID, teamID, roleName string) error
	RemoveRole(ctx context.Context, modelType, modelID, teamID, roleName string) error
	GetRoleNames(ctx context.Context, modelType, modelID, teamID string) ([]string, error)
	GetAllPermissions(ctx context.Context, modelType, modelID, teamID string) ([]string, error)
	GetPermissionsForTeams(ctx context.Context, modelType, modelID string, teamIDs []string) (map[string][]string, error)
}
