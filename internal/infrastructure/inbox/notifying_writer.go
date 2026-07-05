package inbox

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

var _ port.InboxWriter = (*NotifyingWriter)(nil)

// EventPublisher receives inbox change notifications for SSE subscribers.
type EventPublisher interface {
	Publish(workspaceID, eventType string, data any)
}

// PublishFunc adapts a function to EventPublisher.
type PublishFunc func(workspaceID, eventType string, data any)

func (f PublishFunc) Publish(workspaceID, eventType string, data any) {
	if f != nil {
		f(workspaceID, eventType, data)
	}
}

// NotifyingWriter wraps an InboxWriter and publishes events after successful writes.
type NotifyingWriter struct {
	inner port.InboxWriter
	pub   EventPublisher
}

func NewNotifyingWriter(inner port.InboxWriter, pub EventPublisher) *NotifyingWriter {
	return &NotifyingWriter{inner: inner, pub: pub}
}

func (w *NotifyingWriter) WriteEmail(ctx context.Context, workspaceID string, email contracts.Email) (string, error) {
	id, err := w.inner.WriteEmail(ctx, workspaceID, email)
	if err == nil && w.pub != nil {
		w.pub.Publish(workspaceID, "email_received", map[string]string{"id": id})
	}
	return id, err
}

func (w *NotifyingWriter) WriteSMS(ctx context.Context, workspaceID string, sms contracts.SMS) (string, error) {
	id, err := w.inner.WriteSMS(ctx, workspaceID, sms)
	if err == nil && w.pub != nil {
		w.pub.Publish(workspaceID, "sms_received", map[string]string{"id": id})
	}
	return id, err
}

func (w *NotifyingWriter) WritePush(ctx context.Context, workspaceID string, push contracts.PushNotification) (string, error) {
	id, err := w.inner.WritePush(ctx, workspaceID, push)
	if err == nil && w.pub != nil {
		w.pub.Publish(workspaceID, "push_received", map[string]string{"id": id})
	}
	return id, err
}

func (w *NotifyingWriter) WriteChat(ctx context.Context, workspaceID string, chat contracts.ChatMessage) (string, error) {
	id, err := w.inner.WriteChat(ctx, workspaceID, chat)
	if err == nil && w.pub != nil {
		w.pub.Publish(workspaceID, "chat_received", map[string]string{"id": id})
	}
	return id, err
}
