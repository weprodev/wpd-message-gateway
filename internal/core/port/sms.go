package port

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// SMSSender defines the contract for sending SMS messages.
type SMSSender interface {
	Send(ctx context.Context, sms *contracts.SMS) (*contracts.SendResult, error)
	Name() string
}
