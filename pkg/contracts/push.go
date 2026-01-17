package contracts

import "context"

// PushNotification represents a push notification to be sent.
type PushNotification struct {
	DeviceTokens []string          `json:"device_tokens"`
	Title        string            `json:"title"`
	Body         string            `json:"body"`
	Data         map[string]string `json:"data,omitempty"`
	Badge        *int              `json:"badge,omitempty"`
	Sound        string            `json:"sound,omitempty"`
}

// PushSender defines the contract for sending push notifications.
type PushSender interface {
	Send(ctx context.Context, notification *PushNotification) (*SendResult, error)
	Name() string
}
