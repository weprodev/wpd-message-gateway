package contracts

import "context"

// OTP represents a OTP Verification message to be sent.
type OTP struct {
	PhoneNumber			[]string		`json:"phone_number,omitempty"`
	SenderUsername		[]string		`json:"sender_username"`
	Code				[]string		`json:"code,omitempty"`
	CodeLength			int				`json:"code_length"`
}

// OTPSender defines the contract for sending OTP verification messages.
type OTPSender interface {
	Send(ctx context.Context, otp *OTP) (*SendResult, error)
	Name() string
}
