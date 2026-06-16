package postgres

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/weprodev/go-pkg/pgsql"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

type StoredMessageRepository struct {
	client *pgsql.PgClient
}

func NewStoredMessageRepository(client *pgsql.PgClient) port.StoredMessageWriter {
	return &StoredMessageRepository{client: client}
}

func (r *StoredMessageRepository) WriteEmail(ctx context.Context, workspaceID string, email contracts.Email) (string, error) {
	return r.write(ctx, workspaceID, "email", email)
}

func (r *StoredMessageRepository) WriteSMS(ctx context.Context, workspaceID string, sms contracts.SMS) (string, error) {
	return r.write(ctx, workspaceID, "sms", sms)
}

func (r *StoredMessageRepository) WritePush(ctx context.Context, workspaceID string, push contracts.PushNotification) (string, error) {
	return r.write(ctx, workspaceID, "push", push)
}

func (r *StoredMessageRepository) WriteChat(ctx context.Context, workspaceID string, chat contracts.ChatMessage) (string, error) {
	return r.write(ctx, workspaceID, "chat", chat)
}

func (r *StoredMessageRepository) write(ctx context.Context, workspaceID, channelType string, payload any) (string, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	var id string
	err = r.client.GetDB(ctx).QueryRowContext(ctx, `
		INSERT INTO stored_messages (workspace_id, channel_type, payload)
		VALUES ($1, $2, $3)
		RETURNING id
	`, workspaceID, channelType, raw).Scan(&id)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to store message payload",
			"error", err, "workspace_id", workspaceID, "channel", channelType)
	}
	return id, err
}
