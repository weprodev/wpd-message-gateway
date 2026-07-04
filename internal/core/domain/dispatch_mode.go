package domain

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidSettingValue = errors.New("invalid setting value")

const (
	// SettingKeyMessageDispatchMode controls memory capture vs integration dispatch.
	SettingKeyMessageDispatchMode = "message_dispatch_mode"
	// SettingKeyStoreMessageContent controls whether message content is captured in the inbox.
	SettingKeyStoreMessageContent = "store_message_content"
)

// MessageDispatchMode represents where the outbound message is routed.
type MessageDispatchMode string

const (
	// DispatchProvider sends through the connected integration.
	DispatchProvider MessageDispatchMode = "provider"
	// DispatchMemory keeps messages in memory only; external providers are not called.
	DispatchMemory MessageDispatchMode = "memory"
)

// DefaultMessageDispatchMode is used when workspace_settings has no value.
func DefaultMessageDispatchMode() MessageDispatchMode {
	return DispatchMemory
}

// ParseMessageDispatchMode returns the mode if s is a valid dispatch mode string.
func ParseMessageDispatchMode(s string) (MessageDispatchMode, bool) {
	switch MessageDispatchMode(strings.TrimSpace(s)) {
	case DispatchProvider, DispatchMemory:
		return MessageDispatchMode(strings.TrimSpace(s)), true
	default:
		return "", false
	}
}

// MessageDispatchConfig holds dispatch routing and inbox content capture settings.
type MessageDispatchConfig struct {
	Mode                MessageDispatchMode
	StoreMessageContent bool
}

// RoutesViaProvider reports whether outbound traffic should use an integration.
func (c MessageDispatchConfig) RoutesViaProvider() bool {
	return c.Mode == DispatchProvider
}

// ShouldCaptureToInbox reports whether message content should be written to the portal inbox.
func (c MessageDispatchConfig) ShouldCaptureToInbox() bool {
	return c.StoreMessageContent
}

// ResolveMessageDispatchConfig reads canonical workspace setting values.
func ResolveMessageDispatchConfig(modeVal, storeVal string) MessageDispatchConfig {
	config := MessageDispatchConfig{
		Mode:                DefaultMessageDispatchMode(),
		StoreMessageContent: false,
	}
	if mode, ok := ParseMessageDispatchMode(modeVal); ok {
		config.Mode = mode
	}
	if strings.TrimSpace(storeVal) == "true" {
		config.StoreMessageContent = true
	}
	return config
}

// ValidateWorkspaceSettingValue validates known workspace setting keys.
func ValidateWorkspaceSettingValue(key, value string) error {
	switch key {
	case SettingKeyMessageDispatchMode:
		if _, ok := ParseMessageDispatchMode(value); !ok {
			return fmt.Errorf("%w: message_dispatch_mode must be memory or provider", ErrInvalidSettingValue)
		}
	case SettingKeyStoreMessageContent:
		switch strings.TrimSpace(value) {
		case "true", "false":
		default:
			return fmt.Errorf("%w: store_message_content must be true or false", ErrInvalidSettingValue)
		}
	}
	return nil
}
