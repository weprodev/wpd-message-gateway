package port

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// ChatSender defines the contract for sending chat/social media messages.
type ChatSender interface {
	Send(ctx context.Context, message *contracts.ChatMessage) (*contracts.SendResult, error)
	Name() string
}
