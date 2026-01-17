package contracts

import "context"

// PushNotification represents a push notification to be sent
type PushNotification struct {
	DeviceTokens []string          `json:"device_tokens"`   // Target device tokens
	Title        string            `json:"title"`           // Notification title
	Body         string            `json:"body"`            // Notification body
	Data         map[string]string `json:"data,omitempty"`  // Custom data payload
	Badge        *int              `json:"badge,omitempty"` // Badge count (optional)
	Sound        string            `json:"sound,omitempty"` // Notification sound (optional)
}

// PushSender defines the contract for sending push notifications
type PushSender interface {
	// Send sends a push notification and returns the result
	Send(ctx context.Context, notification *PushNotification) (*SendResult, error)

	// Name returns the provider name (e.g., "firebase", "apns")
	Name() string
}
