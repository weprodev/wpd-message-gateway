package domain

import (
	"errors"
	"testing"
)

func TestDefaultMessageDispatchMode(t *testing.T) {
	if got := DefaultMessageDispatchMode(); got != DispatchMemory {
		t.Fatalf("got %q, want %q", got, DispatchMemory)
	}
}

func TestParseMessageDispatchMode(t *testing.T) {
	tests := []struct {
		value string
		want  MessageDispatchMode
		ok    bool
	}{
		{"memory", DispatchMemory, true},
		{"provider", DispatchProvider, true},
		{" memory ", DispatchMemory, true},
		{"unknown", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got, ok := ParseMessageDispatchMode(tt.value)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveMessageDispatchConfig(t *testing.T) {
	tests := []struct {
		name     string
		modeVal  string
		storeVal string
		want     MessageDispatchConfig
	}{
		{"defaults", "", "", MessageDispatchConfig{Mode: DispatchMemory, StoreMessageContent: false}},
		{"memory without store", "memory", "false", MessageDispatchConfig{Mode: DispatchMemory, StoreMessageContent: false}},
		{"provider with store", "provider", "true", MessageDispatchConfig{Mode: DispatchProvider, StoreMessageContent: true}},
		{"invalid mode uses default", "unknown", "true", MessageDispatchConfig{Mode: DispatchMemory, StoreMessageContent: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveMessageDispatchConfig(tt.modeVal, tt.storeVal)
			if got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestMessageDispatchConfig_ShouldCaptureToInbox(t *testing.T) {
	tests := []struct {
		name          string
		config        MessageDispatchConfig
		effectiveMode MessageDispatchMode
		want          bool
	}{
		{"memory always captures", MessageDispatchConfig{Mode: DispatchMemory}, DispatchMemory, true},
		{"provider without store", MessageDispatchConfig{Mode: DispatchProvider}, DispatchProvider, false},
		{"provider with store", MessageDispatchConfig{Mode: DispatchProvider, StoreMessageContent: true}, DispatchProvider, true},
		{"provider fallback to memory", MessageDispatchConfig{Mode: DispatchProvider}, DispatchMemory, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.ShouldCaptureToInbox(tt.effectiveMode); got != tt.want {
				t.Fatalf("ShouldCaptureToInbox() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateWorkspaceSettingValue(t *testing.T) {
	tests := []struct {
		key   string
		value string
		ok    bool
	}{
		{SettingKeyMessageDispatchMode, "memory", true},
		{SettingKeyMessageDispatchMode, "provider", true},
		{SettingKeyMessageDispatchMode, "invalid", false},
		{SettingKeyStoreMessageContent, "true", true},
		{SettingKeyStoreMessageContent, "false", true},
		{SettingKeyStoreMessageContent, "yes", false},
		{"owner_email", "any@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.key+"="+tt.value, func(t *testing.T) {
			err := ValidateWorkspaceSettingValue(tt.key, tt.value)
			if tt.ok && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.ok && !errors.Is(err, ErrInvalidSettingValue) {
				t.Fatalf("expected ErrInvalidSettingValue, got %v", err)
			}
		})
	}
}
