package domain

// SettingKeyDataRetention is the workspace_settings key for portal data retention.
const SettingKeyDataRetention = "data_retention"

// SettingKeyMessageDispatchMode is a legacy key; kept for reading old workspace rows.
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

// DataRetentionValueToDispatchMode maps portal data_retention values to dispatch modes.
func DataRetentionValueToDispatchMode(value string) (MessageDispatchMode, bool) {
	switch value {
	case "memory":
		return DispatchMemoryOnly, true
	case "both":
		return DispatchMemoryAndProvider, true
	case "providers":
		return DispatchProviderOnly, true
	case "provider_database":
		return DispatchProviderAndDatabase, true
	default:
		return ParseMessageDispatchMode(value)
	}
}

// DispatchModeToDataRetentionValue maps dispatch modes to portal data_retention values.
func DispatchModeToDataRetentionValue(value string) (string, bool) {
	if m, ok := ParseMessageDispatchMode(value); ok {
		switch m {
		case DispatchMemoryOnly:
			return "memory", true
		case DispatchMemoryAndProvider:
			return "both", true
		case DispatchProviderOnly:
			return "providers", true
		case DispatchProviderAndDatabase:
			return "provider_database", true
		}
	}
	if _, ok := DataRetentionValueToDispatchMode(value); ok {
		return value, true
	}
	return "", false
}
