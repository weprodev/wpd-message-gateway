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

// OTPStatusChecker is an optional interface that OTPSenders can implement
// to support querying the delivery status of a previously sent message by request_id.
type OTPStatusChecker interface {
	CheckStatus(ctx context.Context, requestID string) (*contracts.VerificationStatus, error)
}
