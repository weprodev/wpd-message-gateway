package contracts

import "context"

// SMS represents an SMS message to be sent.
type SMS struct {
	From    string   `json:"from,omitempty"`
	To      []string `json:"to"`
	Message string   `json:"message"`
}

// SMSSender defines the contract for sending SMS messages.
type SMSSender interface {
	Send(ctx context.Context, sms *SMS) (*SendResult, error)
	Name() string
}
