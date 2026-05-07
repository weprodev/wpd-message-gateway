package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// OTPProvider implements port.OTPSender using an in-memory store.
type OTPProvider struct {
	store *Store
}

// NewOTPProvider creates a new memory OTP provider.
func NewOTPProvider(store *Store) *OTPProvider {
	return &OTPProvider{store: store}
}

// Store returns the underlying memory store.
func (o *OTPProvider) Store() *Store {
	return o.store
}

// Name returns the provider name.
func (o *OTPProvider) Name() string {
	return ProviderName
}

// Send stores the OTP in memory and returns a success result.
func (o *OTPProvider) Send(ctx context.Context, otp *contracts.OTP) (*contracts.SendResult, error) {
	id := uuid.New().String()

	stored := &StoredOTP{
		ID:        id,
		CreatedAt: time.Now(),
		OTP:       otp,
	}
	o.store.AddOTP(stored)

	return &contracts.SendResult{
		ID:         id,
		StatusCode: 200,
		Message:    "Stored OTP in memory",
	}, nil
}