package domain

// SettingKeyDataRetention is the workspace_settings key for portal data retention.
const SettingKeyDataRetention = "data_retention"

// SettingKeyMessageDispatchMode is a legacy key; kept for reading old workspace rows.
const SettingKeyMessageDispatchMode = "message_dispatch_mode"

// Canonical portal data_retention values (aligned with the Portal UI / frontend).
const (
	RetentionMemory           = "memory"
	RetentionMemoryDatabase   = "memory_database"
	RetentionProviders        = "providers"
	RetentionProviderDatabase = "provider_database"
)

// MessageDispatchMode controls memory capture vs integration dispatch for each workspace.
// Provider configs live in DB (integrations); memory is always available as a capture path.
type MessageDispatchMode string

const (
	// DispatchProviderOnly sends only through the connected integration; nothing is kept in process memory.
	DispatchProviderOnly MessageDispatchMode = "provider_only"
	// DispatchProviderAndDatabase sends through the integration and persists the payload in PostgreSQL.
	DispatchProviderAndDatabase MessageDispatchMode = "provider_and_database"
	// DispatchMemoryAndProvider stores in memory and PostgreSQL (stored_messages) and sends through the integration.
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

func retentionValueForMode(mode MessageDispatchMode) (string, bool) {
	switch mode {
	case DispatchMemoryOnly:
		return RetentionMemory, true
	case DispatchMemoryAndProvider:
		return RetentionMemoryDatabase, true
	case DispatchProviderOnly:
		return RetentionProviders, true
	case DispatchProviderAndDatabase:
		return RetentionProviderDatabase, true
	default:
		return "", false
	}
}

// normalizeRetentionValue maps data_retention inputs (including legacy aliases) to canonical portal values.
func normalizeRetentionValue(value string) (string, bool) {
	switch value {
	case RetentionMemory:
		return RetentionMemory, true
	case RetentionProviders, "provider":
		return RetentionProviders, true
	case RetentionMemoryDatabase, "memory_and_database", "both":
		return RetentionMemoryDatabase, true
	case RetentionProviderDatabase, "provider_and_database":
		return RetentionProviderDatabase, true
	default:
		return "", false
	}
}

// DataRetentionValueToDispatchMode maps portal data_retention values to dispatch modes.
func DataRetentionValueToDispatchMode(value string) (MessageDispatchMode, bool) {
	if normalized, ok := normalizeRetentionValue(value); ok {
		value = normalized
	}
	switch value {
	case RetentionMemory:
		return DispatchMemoryOnly, true
	case RetentionMemoryDatabase:
		return DispatchMemoryAndProvider, true
	case RetentionProviders:
		return DispatchProviderOnly, true
	case RetentionProviderDatabase:
		return DispatchProviderAndDatabase, true
	default:
		return ParseMessageDispatchMode(value)
	}
}

// DispatchModeToDataRetentionValue maps dispatch modes or retention values to the canonical portal key.
func DispatchModeToDataRetentionValue(value string) (string, bool) {
	if m, ok := ParseMessageDispatchMode(value); ok {
		return retentionValueForMode(m)
	}
	if normalized, ok := normalizeRetentionValue(value); ok {
		if m, ok := DataRetentionValueToDispatchMode(normalized); ok {
			return retentionValueForMode(m)
		}
	}
	return "", false
}
