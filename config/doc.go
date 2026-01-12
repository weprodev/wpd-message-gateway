// Package config provides configuration loading and management for message providers.
//
// The configuration system supports multiple message provider types (Email, SMS, Push, Chat)
// and automatically discovers provider configurations from environment variables with the
// MESSAGE_ prefix.
//
// # Environment Variable Format
//
// Provider configurations are auto-discovered from environment variables:
//
//	MESSAGE_<PROVIDER>_API_KEY=...
//	MESSAGE_<PROVIDER>_API_SECRET=...
//	MESSAGE_<PROVIDER>_REGION=...
//	MESSAGE_<PROVIDER>_BASE_URL=...
//
// Type-specific fields are also loaded:
//
//	Email: MESSAGE_<PROVIDER>_DOMAIN, MESSAGE_<PROVIDER>_FROM_EMAIL, MESSAGE_<PROVIDER>_FROM_NAME
//	SMS:   MESSAGE_<PROVIDER>_FROM_PHONE
//	Push:  MESSAGE_<PROVIDER>_APP_ID, MESSAGE_<PROVIDER>_TOPIC
//	Chat:  MESSAGE_<PROVIDER>_FROM_PHONE, MESSAGE_<PROVIDER>_WEBHOOK_URL
//
// # Usage
//
//	cfg, err := config.LoadFromEnv()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Access email providers
//	mailgun := cfg.EmailProviders["mailgun"]
//
// # Extending with Custom Providers
//
// Register custom providers before calling LoadFromEnv:
//
//	config.RegisterProvider("myservice", config.ProviderTypeEmail)
//	cfg, _ := config.LoadFromEnv() // Now recognizes MESSAGE_MYSERVICE_* vars
package config
