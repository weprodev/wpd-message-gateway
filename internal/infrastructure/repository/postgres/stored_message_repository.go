package postgres

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/weprodev/go-pkg/pgsql"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
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
		INSERT INTO stored_messages (workspace_id, channel_type, payload, dispatch_status)
		VALUES ($1, $2, $3, 'pending')
		RETURNING id
	`, workspaceID, channelType, raw).Scan(&id)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to store message payload",
			"error", err, "workspace_id", workspaceID, "channel", channelType)
	}
	return id, err
}

func (r *StoredMessageRepository) RecordDispatchOutcome(ctx context.Context, storedMessageID string, outcome domain.StoredMessageDispatchOutcome) error {
	var providerMessageID interface{}
	if outcome.ProviderMessageID != "" {
		providerMessageID = outcome.ProviderMessageID
	}
	var providerStatusCode interface{}
	if outcome.ProviderStatusCode > 0 {
		providerStatusCode = outcome.ProviderStatusCode
	}
	var dispatchError interface{}
	if outcome.DispatchError != "" {
		dispatchError = outcome.DispatchError
	}

	_, err := r.client.GetDB(ctx).ExecContext(ctx, `
		UPDATE stored_messages
		SET dispatch_status = $2,
		    provider_message_id = $3,
		    provider_status_code = $4,
		    dispatch_error = $5,
		    dispatched_at = $6
		WHERE id = $1
	`, storedMessageID, string(outcome.Status), providerMessageID, providerStatusCode, dispatchError, outcome.DispatchedAt)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to record stored message dispatch outcome",
			"error", err, "stored_message_id", storedMessageID, "dispatch_status", outcome.Status)
	}
	return err
}

var _ port.StoredMessageWriter = (*StoredMessageRepository)(nil)
