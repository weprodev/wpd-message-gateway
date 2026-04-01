package port

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

type WorkspaceRepository interface {
	Create(ctx context.Context, workspace *domain.Workspace) error
	GetByID(ctx context.Context, id string) (*domain.Workspace, error)
	GetByUniqueKey(ctx context.Context, uniqueKey string) (*domain.Workspace, error)
	Update(ctx context.Context, workspace *domain.Workspace) error
	SetStatus(ctx context.Context, id, status string) error
	ListForUser(ctx context.Context, userID string) ([]domain.Workspace, error)
}

type APIKeyRepository interface {
	Create(ctx context.Context, apiKey *domain.APIKey) error
	GetByClientID(ctx context.Context, clientID string) (*domain.APIKey, error)
	GetByID(ctx context.Context, id string) (*domain.APIKey, error)
	ListByWorkspace(ctx context.Context, workspaceID string) ([]domain.APIKey, error)
	Delete(ctx context.Context, id string) error
	UpdateSecret(ctx context.Context, id, clientID, secretHash string) error
	UpdateLastUsedAt(ctx context.Context, id string) error
}

type IntegrationRepository interface {
	Create(ctx context.Context, integration *domain.Integration) error
	GetActiveByWorkspaceAndChannel(ctx context.Context, workspaceID, channelType string) (*domain.Integration, error)
	ListByWorkspace(ctx context.Context, workspaceID string) ([]domain.Integration, error)
	GetByID(ctx context.Context, id string) (*domain.Integration, error)
	Delete(ctx context.Context, id string) error
	Upsert(ctx context.Context, integration *domain.Integration) error
}

type TemplateRepository interface {
	Create(ctx context.Context, template *domain.Template) error
	GetByWorkspaceAndKey(ctx context.Context, workspaceID, uniqueKey string) (*domain.Template, error)
	GetByID(ctx context.Context, id string) (*domain.Template, error)
	ListByWorkspace(ctx context.Context, workspaceID string) ([]domain.Template, error)
	Update(ctx context.Context, template *domain.Template) error
	Delete(ctx context.Context, id string) error
}

type MessageRequestLogRepository interface {
	Create(ctx context.Context, log *domain.MessageRequestLog) error
	ListWithSource(ctx context.Context, q MessageLogQuery) ([]domain.MessageRequestLogWithSource, int, error)
}
