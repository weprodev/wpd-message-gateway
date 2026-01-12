package contracts

import "context"

// PushNotification represents a push notification to be sent
type PushNotification struct {
	DeviceTokens []string          // Target device tokens
	Title        string            // Notification title
	Body         string            // Notification body
	Data         map[string]string // Custom data payload
	Badge        *int              // Badge count (optional)
	Sound        string            // Notification sound (optional)
}

// PushSender defines the contract for sending push notifications
type PushSender interface {
	// Send sends a push notification and returns the result
	Send(ctx context.Context, notification *PushNotification) (*SendResult, error)

	// Name returns the provider name (e.g., "firebase", "apns")
	Name() string
}
