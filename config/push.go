package config

// PushConfig specific configuration for Push providers
type PushConfig struct {
	CommonConfig
	AppID string // e.g. OneSignal App ID
	Topic string // e.g. APNS Topic / Bundle ID
}
