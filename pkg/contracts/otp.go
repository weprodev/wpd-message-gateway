package contracts

import "context"

// OTP represents an OTP verification message to be sent.
type OTP struct {
	PhoneNumber    []string `json:"phone_number,omitempty"`
	SenderUsername []string `json:"sender_username,omitempty"`
	CodeLength     int      `json:"code_length"`
}

// OTPSender defines the contract for sending OTP verification messages.
type OTPSender interface {
	Send(ctx context.Context, otp *OTP) (*SendResult, error)
	Name() string
}
