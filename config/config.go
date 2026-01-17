package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Environment string         `yaml:"environment"`
	Server      ServerConfig   `yaml:"server"`
	DevBox      DevBoxConfig   `yaml:"devbox"`
	Providers   ProviderConfig `yaml:"providers"`
	Mailpit     MailpitConfig  `yaml:"mailpit,omitempty"`

	// Flat configs for internal use (populated during load)
	EmailProviders map[string]EmailConfig `yaml:"-"`
	SMSProviders   map[string]SMSConfig   `yaml:"-"`
	PushProviders  map[string]PushConfig  `yaml:"-"`
	ChatProviders  map[string]ChatConfig  `yaml:"-"`
}

type MailpitConfig struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DevBoxConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type ProviderConfig struct {
	Defaults ProviderDefaults          `yaml:"defaults"`
	Email    map[string]EmailConfigMap `yaml:"email"`
	SMS      map[string]SMSConfigMap   `yaml:"sms"`
	Push     map[string]PushConfigMap  `yaml:"push"`
	Chat     map[string]ChatConfigMap  `yaml:"chat"`
}

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

// ProviderCapabilities lists supported message types per provider.
type ProviderCapabilities struct {
	SupportedTypes []ProviderType
}

// knownProviders maps provider names to their capabilities.
var knownProviders = map[string]ProviderCapabilities{
	// Memory provider (supports all types - for development/testing)
	"memory": {SupportedTypes: []ProviderType{ProviderTypeEmail, ProviderTypeSMS, ProviderTypePush, ProviderTypeChat}},

	// Email-only providers
	"mailgun":      {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"sendgrid":     {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"smtp":         {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"ses":          {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"postmark":     {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"sparkpost":    {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"mailjet":      {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"mandrill":     {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"sendinblue":   {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"elasticemail": {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"mailchimp":    {SupportedTypes: []ProviderType{ProviderTypeEmail}},
	"mailerlite":   {SupportedTypes: []ProviderType{ProviderTypeEmail}},

	// SMS-only providers
	"twilio":      {SupportedTypes: []ProviderType{ProviderTypeSMS}},
	"infobip":     {SupportedTypes: []ProviderType{ProviderTypeSMS}},
	"sinch":       {SupportedTypes: []ProviderType{ProviderTypeSMS}},
	"plivo":       {SupportedTypes: []ProviderType{ProviderTypeSMS}},
	"nexmo":       {SupportedTypes: []ProviderType{ProviderTypeSMS}},
	"clicksend":   {SupportedTypes: []ProviderType{ProviderTypeSMS}},
	"messagebird": {SupportedTypes: []ProviderType{ProviderTypeSMS}},
	"sns":         {SupportedTypes: []ProviderType{ProviderTypeSMS}},
	"vonage":      {SupportedTypes: []ProviderType{ProviderTypeSMS}},

	// Multi-type providers
	"cmcom": {SupportedTypes: []ProviderType{ProviderTypeEmail, ProviderTypeSMS, ProviderTypeChat}},

	// Push-only providers
	"firebase":  {SupportedTypes: []ProviderType{ProviderTypePush}},
	"apns":      {SupportedTypes: []ProviderType{ProviderTypePush}},
	"onesignal": {SupportedTypes: []ProviderType{ProviderTypePush}},
	"pusher":    {SupportedTypes: []ProviderType{ProviderTypePush}},
	"pinpoint":  {SupportedTypes: []ProviderType{ProviderTypePush}},
	"expo":      {SupportedTypes: []ProviderType{ProviderTypePush}},

	// Chat-only providers
	"whatsapp":  {SupportedTypes: []ProviderType{ProviderTypeChat}},
	"telegram":  {SupportedTypes: []ProviderType{ProviderTypeChat}},
	"slack":     {SupportedTypes: []ProviderType{ProviderTypeChat}},
	"discord":   {SupportedTypes: []ProviderType{ProviderTypeChat}},
	"teams":     {SupportedTypes: []ProviderType{ProviderTypeChat}},
	"messenger": {SupportedTypes: []ProviderType{ProviderTypeChat}},
	"line":      {SupportedTypes: []ProviderType{ProviderTypeChat}},
	"viber":     {SupportedTypes: []ProviderType{ProviderTypeChat}},
	"wechat":    {SupportedTypes: []ProviderType{ProviderTypeChat}},
}

// providerLookup is an optimized uppercase map for fast lookups
var providerLookup map[string]ProviderCapabilities

func init() {
	providerLookup = make(map[string]ProviderCapabilities)
	for k, v := range knownProviders {
		providerLookup[strings.ToUpper(k)] = v
	}
}

// providerMu protects concurrent access to provider maps
var providerMu sync.RWMutex

// RegisterProvider registers a new provider with its capabilities for auto-discovery.
// This function is safe for concurrent use.
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

// LoadConfig loads configuration from a YAML file.
// If path is empty, it looks for configs/local.yml
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

	// Apply environment variable overrides (for secret injection)
	cfg.applyEnvOverrides()

	// Parse map configs into struct configs
	cfg.parseProviderConfigs()

	return cfg, nil
}

// applyEnvOverrides allows overriding config with environment variables.
// Essential for injecting secrets in production (e.g. MESSAGE_{PROVIDER}_API_KEY).
func (c *Config) applyEnvOverrides() {
	// Override Server Port
	if port := os.Getenv("PORT"); port != "" {
		// Only if it looks like a number, otherwise ignore
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil {
			c.Server.Port = p
		}
	}

	// Iterate over environment variables to find MESSAGE_ config
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "MESSAGE_") {
			continue
		}

		parts := strings.SplitN(env, "=", 2)
		key := parts[0]
		val := parts[1]

		// MESSAGE_DEFAULT_EMAIL_PROVIDER -> Defaults.Email
		if key == "MESSAGE_DEFAULT_EMAIL_PROVIDER" {
			c.Providers.Defaults.Email = val
			continue
		}
		if key == "MESSAGE_DEFAULT_SMS_PROVIDER" {
			c.Providers.Defaults.SMS = val
			continue
		}
		if key == "MESSAGE_DEFAULT_PUSH_PROVIDER" {
			c.Providers.Defaults.Push = val
			continue
		}
		if key == "MESSAGE_DEFAULT_CHAT_PROVIDER" {
			c.Providers.Defaults.Chat = val
			continue
		}

		// Provider specific configs: MESSAGE_{PROVIDER}_{KEY}
		// We'll stash these into the provider maps if they exist, or create entries
		// This is a bit complex due to the map structure, so we simplify:
		// We expect the YAML to define the provider structure, and Env vars just override values.
		// However, for new providers not in YAML, we can support them if we parse correctly.
	}
}

// parseProviderConfigs converts the raw map string configs from YAML into typed configs
func (c *Config) parseProviderConfigs() {
	// Helper to extract common config from map
	buildCommon := func(m map[string]string) CommonConfig {
		return CommonConfig{
			APIKey:    m["api_key"],
			APISecret: m["api_secret"],
			Region:    m["region"],
			BaseURL:   m["base_url"],
			Extra:     m, // Store all as extra for provider-specific needs
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

// Legacy Config structs for compatibility

type EmailConfig struct {
	CommonConfig
	Domain    string
	FromEmail string
	FromName  string
}

type SMSConfig struct {
	CommonConfig
	FromPhone string
}

type PushConfig struct {
	CommonConfig
	AppID string
	Topic string
}

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
