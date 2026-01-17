package app

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration.
type Config struct {
	Environment string         `yaml:"environment"`
	Server      ServerConfig   `yaml:"server"`
	DevBox      DevBoxConfig   `yaml:"devbox"`
	Providers   ProviderConfig `yaml:"providers"`
	Mailpit     MailpitConfig  `yaml:"mailpit,omitempty"`

	// Parsed provider configs
	EmailProviders map[string]EmailConfig `yaml:"-"`
	SMSProviders   map[string]SMSConfig   `yaml:"-"`
	PushProviders  map[string]PushConfig  `yaml:"-"`
	ChatProviders  map[string]ChatConfig  `yaml:"-"`
}

// MailpitConfig holds SMTP forwarding configuration.
type MailpitConfig struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

// ServerConfig holds server configuration.
type ServerConfig struct {
	Port int `yaml:"port"`
}

// DevBoxConfig holds devbox configuration.
type DevBoxConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

// ProviderConfig holds provider configuration.
type ProviderConfig struct {
	Defaults ProviderDefaults          `yaml:"defaults"`
	Email    map[string]EmailConfigMap `yaml:"email"`
	SMS      map[string]SMSConfigMap   `yaml:"sms"`
	Push     map[string]PushConfigMap  `yaml:"push"`
	Chat     map[string]ChatConfigMap  `yaml:"chat"`
}

// ProviderDefaults holds default provider names.
type ProviderDefaults struct {
	Email string `yaml:"email"`
	SMS   string `yaml:"sms"`
	Push  string `yaml:"push"`
	Chat  string `yaml:"chat"`
}

// Intermediate maps for YAML parsing
type EmailConfigMap map[string]string
type SMSConfigMap map[string]string
type PushConfigMap map[string]string
type ChatConfigMap map[string]string

// CommonConfig shared fields across all providers.
type CommonConfig struct {
	APIKey    string
	APISecret string
	Region    string
	BaseURL   string
	Extra     map[string]string
}

// EmailConfig holds email provider configuration.
type EmailConfig struct {
	CommonConfig
	Domain    string
	FromEmail string
	FromName  string
}

// SMSConfig holds SMS provider configuration.
type SMSConfig struct {
	CommonConfig
	FromPhone string
}

// PushConfig holds push notification provider configuration.
type PushConfig struct {
	CommonConfig
	AppID string
	Topic string
}

// ChatConfig holds chat provider configuration.
type ChatConfig struct {
	CommonConfig
	FromPhone  string
	WebhookURL string
}

// Default Providers Helpers
func (c *Config) DefaultEmailProvider() string { return c.Providers.Defaults.Email }
func (c *Config) DefaultSMSProvider() string   { return c.Providers.Defaults.SMS }
func (c *Config) DefaultPushProvider() string  { return c.Providers.Defaults.Push }
func (c *Config) DefaultChatProvider() string  { return c.Providers.Defaults.Chat }

// LoadConfig loads configuration from a YAML file.
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = "configs/local.yml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	cfg := &Config{
		EmailProviders: make(map[string]EmailConfig),
		SMSProviders:   make(map[string]SMSConfig),
		PushProviders:  make(map[string]PushConfig),
		ChatProviders:  make(map[string]ChatConfig),
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	cfg.applyEnvOverrides()
	cfg.parseProviderConfigs()

	return cfg, nil
}

// MustLoadConfig loads config and panics on error.
func MustLoadConfig(path string) *Config {
	cfg, err := LoadConfig(path)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

// applyEnvOverrides applies environment variable overrides.
func (c *Config) applyEnvOverrides() {
	if port := os.Getenv("PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil {
			c.Server.Port = p
		}
	}

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "MESSAGE_") {
			continue
		}

		parts := strings.SplitN(env, "=", 2)
		key := parts[0]
		val := parts[1]

		switch key {
		case "MESSAGE_DEFAULT_EMAIL_PROVIDER":
			c.Providers.Defaults.Email = val
		case "MESSAGE_DEFAULT_SMS_PROVIDER":
			c.Providers.Defaults.SMS = val
		case "MESSAGE_DEFAULT_PUSH_PROVIDER":
			c.Providers.Defaults.Push = val
		case "MESSAGE_DEFAULT_CHAT_PROVIDER":
			c.Providers.Defaults.Chat = val
		}
	}
}

// parseProviderConfigs converts raw map configs into typed configs.
func (c *Config) parseProviderConfigs() {
	buildCommon := func(m map[string]string) CommonConfig {
		return CommonConfig{
			APIKey:    m["api_key"],
			APISecret: m["api_secret"],
			Region:    m["region"],
			BaseURL:   m["base_url"],
			Extra:     m,
		}
	}

	for name, m := range c.Providers.Email {
		c.EmailProviders[name] = EmailConfig{
			CommonConfig: buildCommon(map[string]string(m)),
			Domain:       m["domain"],
			FromEmail:    m["from_email"],
			FromName:     m["from_name"],
		}
	}

	for name, m := range c.Providers.SMS {
		c.SMSProviders[name] = SMSConfig{
			CommonConfig: buildCommon(map[string]string(m)),
			FromPhone:    m["from_phone"],
		}
	}

	for name, m := range c.Providers.Push {
		c.PushProviders[name] = PushConfig{
			CommonConfig: buildCommon(map[string]string(m)),
			AppID:        m["app_id"],
			Topic:        m["topic"],
		}
	}

	for name, m := range c.Providers.Chat {
		c.ChatProviders[name] = ChatConfig{
			CommonConfig: buildCommon(map[string]string(m)),
			FromPhone:    m["from_phone"],
			WebhookURL:   m["webhook_url"],
		}
	}
}

// --- Provider Type and Capabilities ---

// ProviderType defines the type of message provider.
type ProviderType string

const (
	ProviderTypeEmail ProviderType = "email"
	ProviderTypeSMS   ProviderType = "sms"
	ProviderTypePush  ProviderType = "push"
	ProviderTypeChat  ProviderType = "chat"
)

// ProviderCapabilities lists supported message types per provider.
type ProviderCapabilities struct {
	SupportedTypes []ProviderType
}

var knownProviders = map[string]ProviderCapabilities{
	"memory":   {SupportedTypes: []ProviderType{ProviderTypeEmail, ProviderTypeSMS, ProviderTypePush, ProviderTypeChat}},
	"mailgun":  {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"sendgrid": {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"smtp":     {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"twilio":   {SupportedTypes: []ProviderType{ProviderTypeSMS}},
	"firebase": {SupportedTypes: []ProviderType{ProviderTypePush}},
	"whatsapp": {SupportedTypes: []ProviderType{ProviderTypeChat}},
	"telegram": {SupportedTypes: []ProviderType{ProviderTypeChat}},
}

var providerLookup map[string]ProviderCapabilities
var providerMu sync.RWMutex

func init() {
	providerLookup = make(map[string]ProviderCapabilities)
	for k, v := range knownProviders {
		providerLookup[strings.ToUpper(k)] = v
	}
}

// RegisterProvider registers a new provider.
func RegisterProvider(name string, capabilities ProviderCapabilities) {
	providerMu.Lock()
	defer providerMu.Unlock()

	lower := strings.ToLower(name)
	knownProviders[lower] = capabilities
	providerLookup[strings.ToUpper(name)] = capabilities
}

// IsKnownProvider checks if a provider name is registered.
func IsKnownProvider(name string) bool {
	providerMu.RLock()
	defer providerMu.RUnlock()
	_, ok := providerLookup[strings.ToUpper(name)]
	return ok
}
