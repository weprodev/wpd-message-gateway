package contracts

import "context"

// SMS represents an SMS message to be sent
type SMS struct {
	From    string   // Sender phone number or alphanumeric ID
	To      []string // Recipient phone numbers (E.164 format recommended)
	Message string   // SMS content (max 160 chars for single SMS)
}

// SMSSender defines the contract for sending SMS messages
type SMSSender interface {
	// Send sends an SMS and returns the result
	Send(ctx context.Context, sms *SMS) (*SendResult, error)

	// Name returns the provider name (e.g., "twilio", "cmcom")
	Name() string
}
