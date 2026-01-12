package config

// EmailConfig specific configuration for Email providers
type EmailConfig struct {
	CommonConfig
	Domain    string
	FromEmail string
	FromName  string
}
