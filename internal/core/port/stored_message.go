package port

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// StoredMessageWriter persists outbound message payloads for durable retention.
type StoredMessageWriter interface {
	WriteEmail(ctx context.Context, workspaceID string, email contracts.Email) (id string, err error)
	WriteSMS(ctx context.Context, workspaceID string, sms contracts.SMS) (id string, err error)
	WritePush(ctx context.Context, workspaceID string, push contracts.PushNotification) (id string, err error)
	WriteChat(ctx context.Context, workspaceID string, chat contracts.ChatMessage) (id string, err error)
}
