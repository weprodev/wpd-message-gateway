package domain

// Workspace setting key for how outbound messages are handled relative to memory vs providers.
const SettingKeyMessageDispatchMode = "message_dispatch_mode"

// MessageDispatchMode controls memory capture vs integration dispatch for each workspace.
// Provider configs live in DB (integrations); memory is always available as a capture path.
type MessageDispatchMode string

const (
	// DispatchProviderOnly sends only through the connected integration; nothing is kept in process memory.
	DispatchProviderOnly MessageDispatchMode = "provider_only"
	// DispatchProviderAndDatabase sends through the integration like provider_only; request logs are retained.
	DispatchProviderAndDatabase MessageDispatchMode = "provider_and_database"
	// DispatchMemoryAndProvider stores in memory and sends through the integration.
	DispatchMemoryAndProvider MessageDispatchMode = "memory_and_provider"
	// DispatchMemoryOnly keeps messages in memory only; external providers are not called.
	DispatchMemoryOnly MessageDispatchMode = "memory_only"
)

// DefaultMessageDispatchMode is used when workspace_settings has no value (matches portal “safe dev” default).
func DefaultMessageDispatchMode() MessageDispatchMode {
	return DispatchMemoryOnly
}

// ParseMessageDispatchMode returns the mode if s is a gateway dispatch string.
func ParseMessageDispatchMode(s string) (MessageDispatchMode, bool) {
	switch MessageDispatchMode(s) {
	case DispatchProviderOnly, DispatchProviderAndDatabase, DispatchMemoryAndProvider, DispatchMemoryOnly:
		return MessageDispatchMode(s), true
	default:
		return "", false
	}
}

// SettingValueToDispatchMode maps portal/API setting values (and gateway strings) to MessageDispatchMode.
func SettingValueToDispatchMode(value string) (MessageDispatchMode, bool) {
	if mode, ok := ParseMessageDispatchMode(value); ok {
		return mode, true
	}
	switch value {
	case "memory":
		return DispatchMemoryOnly, true
	case "both", "memory_database":
		return DispatchMemoryAndProvider, true
	case "providers", "provider":
		return DispatchProviderOnly, true
	case "provider_database":
		return DispatchProviderAndDatabase, true
	default:
		return "", false
	}
}

// ShouldRetainRequestLog reports whether request logs for this dispatch mode are long-term retained.
func ShouldRetainRequestLog(mode MessageDispatchMode) bool {
	switch mode {
	case DispatchMemoryAndProvider, DispatchProviderAndDatabase:
		return true
	default:
		return false
	}
}
