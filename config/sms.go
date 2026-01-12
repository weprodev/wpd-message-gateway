package config

// SMSConfig specific configuration for SMS providers
type SMSConfig struct {
	CommonConfig
	FromPhone string
}
