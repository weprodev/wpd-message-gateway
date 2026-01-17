package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/contracts"
)

// ChatProvider wraps the Store to implement contracts.ChatSender
type ChatProvider struct {
	store *Provider
}

// Store returns the underlying memory store
func (c *ChatProvider) Store() *Provider {
	return c.store
}

// Name returns the provider name.
func (c *ChatProvider) Name() string {
	return ProviderName
}

// Send stores the chat message in memory and returns a success result.
func (c *ChatProvider) Send(ctx context.Context, chat *contracts.ChatMessage) (*contracts.SendResult, error) {
	id := uuid.New().String()

	stored := &StoredChat{
		ID:        id,
		CreatedAt: time.Now(),
		Chat:      chat,
	}
	c.store.addChat(stored)

	return &contracts.SendResult{
		ID:         id,
		StatusCode: 200,
		Message:    "Stored chat message in memory",
	}, nil
}
