package port

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// OTPSender defines the contract for sending OTP Verification messages.
type OTPSender interface {
	Send(ctx context.Context, otp *contracts.OTP) (*contracts.SendResult, error)
	Name() string
}
