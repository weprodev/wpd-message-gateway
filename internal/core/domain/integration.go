package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Integration lifecycle statuses (integrations.status column).
const (
	IntegrationStatusConnected    = "connected"
	IntegrationStatusDisconnected = "disconnected"
)

// ProviderNameMemory is the in-process capture provider; not valid for provider dispatch mode.
const ProviderNameMemory = "memory"

// ErrProviderNotReady indicates the workspace integration cannot send via an external provider.
var ErrProviderNotReady = errors.New("provider not ready")

// Integration holds provider credentials for a workspace channel.
type Integration struct {
	ID           string    `json:"id"`
	WorkspaceID  string    `json:"workspace_id"`
	ChannelType  string    `json:"channel_type"`
	ProviderName string    `json:"provider_name"`
	Config       []byte    `json:"-"` // decrypted JSON
	Status       string    `json:"status"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ValidateProviderIntegration checks that an integration can be used for provider dispatch.
func ValidateProviderIntegration(intg Integration) error {
	if intg.ProviderName == ProviderNameMemory {
		return fmt.Errorf("%w: memory provider is not allowed in provider dispatch mode", ErrProviderNotReady)
	}
	if intg.Status != IntegrationStatusConnected {
		return fmt.Errorf("%w: integration status is %q", ErrProviderNotReady, intg.Status)
	}
	if len(strings.TrimSpace(intg.ProviderName)) == 0 {
		return fmt.Errorf("%w: provider name is empty", ErrProviderNotReady)
	}
	if len(strings.TrimSpace(string(intg.Config))) == 0 {
		return fmt.Errorf("%w: provider config is empty", ErrProviderNotReady)
	}
	if !json.Valid(intg.Config) {
		return fmt.Errorf("%w: provider config is invalid JSON", ErrProviderNotReady)
	}
	return nil
}
