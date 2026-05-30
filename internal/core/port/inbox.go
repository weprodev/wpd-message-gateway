package port

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// InboxWriter is the write side of the in-process message capture store.
//
// Implementations live in the infrastructure layer (memory provider).
// The Core domain depends only on this interface — never on the concrete store.
// Callers receive a message ID that can be stored alongside the provider result.
type InboxWriter interface {
	// WriteEmail captures an email against workspaceID and returns an assigned ID.
	WriteEmail(ctx context.Context, workspaceID string, email *contracts.Email) (id string, err error)

	// WriteSMS captures an SMS against workspaceID and returns an assigned ID.
	WriteSMS(ctx context.Context, workspaceID string, sms *contracts.SMS) (id string, err error)

	// WritePush captures a push notification against workspaceID and returns an assigned ID.
	WritePush(ctx context.Context, workspaceID string, push *contracts.PushNotification) (id string, err error)

	// WriteChat captures a chat message against workspaceID and returns an assigned ID.
	WriteChat(ctx context.Context, workspaceID string, chat *contracts.ChatMessage) (id string, err error)

	// WriteOTP captures a otp verification message against workspaceID and returns an assigned ID.
	WriteOTP(ctx context.Context, workspaceID string, otp *contracts.OTP) (id string, err error)
}
