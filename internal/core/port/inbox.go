package port

import (
	"context"
	"time"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// StoredEmail wraps an email with metadata for storage.
type StoredEmail struct {
	ID          string          `json:"id"`
	WorkspaceID string          `json:"workspace_id,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	Email       contracts.Email `json:"email"`
}

// StoredSMS wraps an SMS with metadata for storage.
type StoredSMS struct {
	ID          string        `json:"id"`
	WorkspaceID string        `json:"workspace_id,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	SMS         contracts.SMS `json:"sms"`
}

// StoredPush wraps a push notification with metadata for storage.
type StoredPush struct {
	ID          string                     `json:"id"`
	WorkspaceID string                     `json:"workspace_id,omitempty"`
	CreatedAt   time.Time                  `json:"created_at"`
	Push        contracts.PushNotification `json:"push"`
}

// StoredChat wraps a chat message with metadata for storage.
type StoredChat struct {
	ID          string                `json:"id"`
	WorkspaceID string                `json:"workspace_id,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	Chat        contracts.ChatMessage `json:"chat"`
}

// InboxReader is the read side of the in-process message capture store.
type InboxReader interface {
	StatsForWorkspace(workspaceID string) map[string]int
	EmailsForWorkspace(workspaceID string) []StoredEmail
	EmailByIDForWorkspace(id, workspaceID string) (StoredEmail, bool)
	DeleteEmailByIDForWorkspace(id, workspaceID string) bool
	SMSForWorkspace(workspaceID string) []StoredSMS
	SMSByIDForWorkspace(id, workspaceID string) (StoredSMS, bool)
	DeleteSMSByIDForWorkspace(id, workspaceID string) bool
	PushForWorkspace(workspaceID string) []StoredPush
	PushByIDForWorkspace(id, workspaceID string) (StoredPush, bool)
	DeletePushByIDForWorkspace(id, workspaceID string) bool
	ChatForWorkspace(workspaceID string) []StoredChat
	ChatByIDForWorkspace(id, workspaceID string) (StoredChat, bool)
	DeleteChatByIDForWorkspace(id, workspaceID string) bool
	ClearWorkspace(workspaceID string)
}

// InboxWriter is the write side of the in-process message capture store.
type InboxWriter interface {
	WriteEmail(ctx context.Context, workspaceID string, email contracts.Email) (id string, err error)
	WriteSMS(ctx context.Context, workspaceID string, sms contracts.SMS) (id string, err error)
	WritePush(ctx context.Context, workspaceID string, push contracts.PushNotification) (id string, err error)
	WriteChat(ctx context.Context, workspaceID string, chat contracts.ChatMessage) (id string, err error)
}
