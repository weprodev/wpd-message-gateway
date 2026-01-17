package port

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// PushSender defines the contract for sending push notifications.
type PushSender interface {
	Send(ctx context.Context, notification *contracts.PushNotification) (*contracts.SendResult, error)
	Name() string
}
