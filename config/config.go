package config

import (
	"os"
	"strings"
	"sync"
)

// CommonConfig shared fields across all providers
type CommonConfig struct {
	APIKey    string
	APISecret string
	Region    string
	BaseURL   string
	Extra     map[string]string
}

// ProviderType defines the type of message provider
type ProviderType string

const (
	ProviderTypeEmail ProviderType = "email"
	ProviderTypeSMS   ProviderType = "sms"
	ProviderTypePush  ProviderType = "push"
	ProviderTypeChat  ProviderType = "chat"
)

// knownProviders maps provider names to their types for auto-discovery
var knownProviders = map[string]ProviderType{
	// Email Providers
	"mailgun":      ProviderTypeEmail,
	"sendgrid":     ProviderTypeEmail,
	"smtp":         ProviderTypeEmail,
	"ses":          ProviderTypeEmail,
	"postmark":     ProviderTypeEmail,
	"sparkpost":    ProviderTypeEmail,
	"mailjet":      ProviderTypeEmail,
	"mandrill":     ProviderTypeEmail,
	"sendinblue":   ProviderTypeEmail,
	"elasticemail": ProviderTypeEmail,
	"mailchimp":    ProviderTypeEmail,
	"mailerlite":   ProviderTypeEmail,

	// SMS Providers
	"twilio":      ProviderTypeSMS,
	"cmcom":       ProviderTypeSMS,
	"infobip":     ProviderTypeSMS,
	"sinch":       ProviderTypeSMS,
	"plivo":       ProviderTypeSMS,
	"nexmo":       ProviderTypeSMS,
	"clicksend":   ProviderTypeSMS,
	"messagebird": ProviderTypeSMS,
	"sns":         ProviderTypeSMS,
	"vonage":      ProviderTypeSMS,

	// Push Providers
	"firebase":  ProviderTypePush,
	"apns":      ProviderTypePush,
	"onesignal": ProviderTypePush,
	"pusher":    ProviderTypePush,
	"pinpoint":  ProviderTypePush,
	"expo":      ProviderTypePush,

	// Chat Providers
	"whatsapp":  ProviderTypeChat,
	"telegram":  ProviderTypeChat,
	"slack":     ProviderTypeChat,
	"discord":   ProviderTypeChat,
	"teams":     ProviderTypeChat,
	"messenger": ProviderTypeChat,
	"line":      ProviderTypeChat,
	"viber":     ProviderTypeChat,
	"wechat":    ProviderTypeChat,
}

// providerLookup is an optimized uppercase map for fast lookups
var providerLookup map[string]ProviderType

func init() {
	providerLookup = make(map[string]ProviderType)
	for k, v := range knownProviders {
		providerLookup[strings.ToUpper(k)] = v
	}
}

// providerMu protects concurrent access to provider maps
var providerMu sync.RWMutex

// RegisterProvider registers a new provider type for auto-discovery.
// This function is safe for concurrent use.
func RegisterProvider(name string, pType ProviderType) {
	providerMu.Lock()
	defer providerMu.Unlock()

	lower := strings.ToLower(name)
	knownProviders[lower] = pType
	providerLookup[strings.ToUpper(name)] = pType
}

// Config holds all provider configurations
type Config struct {
	DefaultEmailProvider string
	DefaultSMSProvider   string
	DefaultPushProvider  string
	DefaultChatProvider  string

	// Typed provider configurations
	EmailProviders map[string]EmailConfig
	SMSProviders   map[string]SMSConfig
	PushProviders  map[string]PushConfig
	ChatProviders  map[string]ChatConfig
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	cfg := NewConfig()

	// Load defaults
	cfg.DefaultEmailProvider = getEnv("MESSAGE_DEFAULT_EMAIL_PROVIDER", "")
	cfg.DefaultSMSProvider = getEnv("MESSAGE_DEFAULT_SMS_PROVIDER", "")
	cfg.DefaultPushProvider = getEnv("MESSAGE_DEFAULT_PUSH_PROVIDER", "")
	cfg.DefaultChatProvider = getEnv("MESSAGE_DEFAULT_CHAT_PROVIDER", "")

	loaded := make(map[string]bool)

	// Discover and load provider configs
	for _, env := range os.Environ() {
		// Optimization: Check prefix "MESSAGE_" (len 8)
		// Skip "MESSAGE_DEFAULT_" (len 16)
		if len(env) <= 8 || env[:8] != "MESSAGE_" || (len(env) > 16 && env[:16] == "MESSAGE_DEFAULT_") {
			continue
		}

		// Find '=' to isolate key
		// e.g. "MESSAGE_MAILGUN_API_KEY=..."
		eq := strings.IndexByte(env, '=')
		if eq == -1 {
			continue
		}

		key := env[:eq] // No allocation, string slicing

		// Extract provider name by iterating segments between "MESSAGE_" and potential suffixes
		// Start after "MESSAGE_" (index 8)
		start := 8

		// Loop through underscores to find potential provider name matches
		// e.g. "MESSAGE_MAILGUN_API_KEY" -> "MAILGUN"
		// e.g. "MESSAGE_CM_MMC_API_KEY" -> "CM" (fail) -> "CM_MMC" (hit)
		for i := start; i < len(key); i++ {
			if key[i] == '_' {
				// Potential name is key[start:i] (e.g. "MAILGUN")
				candidateUpper := key[start:i]

				// O(1) Lookup in optimized map
				if pType, known := providerLookup[candidateUpper]; known {
					// Found!
					// Convert to lowercase for canonical storage (only allocs on success)
					candidateLower := strings.ToLower(candidateUpper)

					if !loaded[candidateLower] {
						cfg.loadProvider(candidateLower, pType)
						loaded[candidateLower] = true
					}
					// Once found, we don't need to check further segments for this key
					break
				}
			}
		}
	}

	return cfg, nil
}

// NewConfig creates a new empty Config
func NewConfig() *Config {
	return &Config{
		EmailProviders: make(map[string]EmailConfig),
		SMSProviders:   make(map[string]SMSConfig),
		PushProviders:  make(map[string]PushConfig),
		ChatProviders:  make(map[string]ChatConfig),
	}
}

func (c *Config) loadProvider(name string, pType ProviderType) {
	prefix := "MESSAGE_" + strings.ToUpper(name) + "_"

	common := CommonConfig{
		APIKey:    getEnv(prefix+"API_KEY", ""),
		APISecret: getEnv(prefix+"API_SECRET", ""),
		Region:    getEnv(prefix+"REGION", ""),
		BaseURL:   getEnv(prefix+"BASE_URL", ""),
		Extra:     make(map[string]string),
	}

	switch pType {
	case ProviderTypeEmail:
		c.EmailProviders[name] = EmailConfig{
			CommonConfig: common,
			Domain:       getEnv(prefix+"DOMAIN", ""),
			FromEmail:    getEnv(prefix+"FROM_EMAIL", ""),
			FromName:     getEnv(prefix+"FROM_NAME", ""),
		}
	case ProviderTypeSMS:
		c.SMSProviders[name] = SMSConfig{
			CommonConfig: common,
			FromPhone:    getEnv(prefix+"FROM_PHONE", ""),
		}
	case ProviderTypePush:
		c.PushProviders[name] = PushConfig{
			CommonConfig: common,
			AppID:        getEnv(prefix+"APP_ID", ""),
			Topic:        getEnv(prefix+"TOPIC", ""),
		}
	case ProviderTypeChat:
		c.ChatProviders[name] = ChatConfig{
			CommonConfig: common,
			FromPhone:    getEnv(prefix+"FROM_PHONE", ""),
			WebhookURL:   getEnv(prefix+"WEBHOOK_URL", ""),
		}
	}
}

// getEnv retrieves an environment variable or returns defaultValue if not set.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
