package domain

// Workspace setting key for how outbound messages are handled relative to memory vs providers.
const SettingKeyMessageDispatchMode = "message_dispatch_mode"

// MessageDispatchMode controls memory capture vs integration dispatch for each workspace.
// Provider configs live in DB (integrations); memory is always available as a capture path.
type MessageDispatchMode string

const (
	// DispatchMemory captures messages in process memory (portal inbox); external providers are not called.
	DispatchMemory MessageDispatchMode = "memory"
	// DispatchProvider sends through the connected integration.
	DispatchProvider MessageDispatchMode = "provider"
)

// MessageDispatchAPIValue is the canonical gateway string stored in settings and stamped in response meta.
type MessageDispatchAPIValue string

const (
	APIMemoryOnly          MessageDispatchAPIValue = "memory_only"
	APIMemoryAndDatabase   MessageDispatchAPIValue = "memory_and_database"
	APIProviderOnly        MessageDispatchAPIValue = "provider_only"
	APIProviderAndDatabase MessageDispatchAPIValue = "provider_and_database"
)

// MessageDispatchConfig pairs dispatch path with request-log retention.
type MessageDispatchConfig struct {
	Mode             MessageDispatchMode
	RetainRequestLog bool
}

// DefaultMessageDispatchConfig is used when workspace_settings has no value (matches portal “safe dev” default).
func DefaultMessageDispatchConfig() MessageDispatchConfig {
	return MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: false}
}

// APIValue returns the canonical gateway string for this config (settings + response meta).
func (c MessageDispatchConfig) APIValue() MessageDispatchAPIValue {
	if c.Mode == DispatchMemory {
		if c.RetainRequestLog {
			return APIMemoryAndDatabase
		}
		return APIMemoryOnly
	}
	if c.RetainRequestLog {
		return APIProviderAndDatabase
	}
	return APIProviderOnly
}

// ParseMessageDispatchConfig parses a gateway API string or legacy setting alias.
func ParseMessageDispatchConfig(s string) (MessageDispatchConfig, bool) {
	switch MessageDispatchAPIValue(s) {
	case APIMemoryOnly:
		return MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: false}, true
	case APIMemoryAndDatabase:
		return MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: true}, true
	case APIProviderOnly:
		return MessageDispatchConfig{Mode: DispatchProvider, RetainRequestLog: false}, true
	case APIProviderAndDatabase:
		return MessageDispatchConfig{Mode: DispatchProvider, RetainRequestLog: true}, true
	}
	return legacySettingToDispatchConfig(s)
}

// SettingValueToDispatchConfig maps workspace_settings.message_dispatch_mode to MessageDispatchConfig.
// Canonical stored values are the four gateway strings (memory_only, memory_and_database, …).
// Short legacy aliases (memory, memory_database, both, provider, …) are accepted on read only.
func SettingValueToDispatchConfig(value string) (MessageDispatchConfig, bool) {
	return ParseMessageDispatchConfig(value)
}

func legacySettingToDispatchConfig(value string) (MessageDispatchConfig, bool) {
	switch value {
	case "memory":
		return MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: false}, true
	case "both", "memory_database", "memory_and_provider":
		return MessageDispatchConfig{Mode: DispatchMemory, RetainRequestLog: true}, true
	case "providers", "provider":
		return MessageDispatchConfig{Mode: DispatchProvider, RetainRequestLog: false}, true
	case "provider_database":
		return MessageDispatchConfig{Mode: DispatchProvider, RetainRequestLog: true}, true
	default:
		return MessageDispatchConfig{}, false
	}
}
