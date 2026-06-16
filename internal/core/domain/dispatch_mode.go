package domain

// Workspace setting key for how outbound messages are handled relative to memory vs providers.
const SettingKeyMessageDispatchMode = "message_dispatch_mode"

// MessageDispatchMode controls memory capture vs integration dispatch for each workspace.
// Provider configs live in DB (integrations); memory is always available as a capture path.
type MessageDispatchMode string

const (
	// DispatchProviderOnly sends only through the connected integration; nothing is kept in process memory.
	DispatchProviderOnly MessageDispatchMode = "provider_only"
	// DispatchProviderAndDatabase sends through the integration and persists the payload in PostgreSQL.
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

// ParseMessageDispatchMode returns the mode if s is valid.
func ParseMessageDispatchMode(s string) (MessageDispatchMode, bool) {
	switch MessageDispatchMode(s) {
	case DispatchProviderOnly, DispatchProviderAndDatabase, DispatchMemoryAndProvider, DispatchMemoryOnly:
		return MessageDispatchMode(s), true
	default:
		return "", false
	}
}
