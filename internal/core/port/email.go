package port

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// EmailSender defines the contract for sending emails.
type EmailSender interface {
	Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error)
	Name() string
}
