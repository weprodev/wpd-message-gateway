package app

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/weprodev/wpd-message-gateway/internal/app/registry"
)

// Config represents the application configuration.
type Config struct {
	Environment string         `yaml:"environment"`
	Server      ServerConfig   `yaml:"server"`
	DevBox      DevBoxConfig   `yaml:"devbox"`
	Providers   ProviderConfig `yaml:"providers"`
	Mailpit     MailpitConfig  `yaml:"mailpit,omitempty"`

	// Parsed provider configs - using registry types as single source of truth
	EmailProviders map[string]registry.EmailConfig `yaml:"-"`
	SMSProviders   map[string]registry.SMSConfig   `yaml:"-"`
	PushProviders  map[string]registry.PushConfig  `yaml:"-"`
	ChatProviders  map[string]registry.ChatConfig  `yaml:"-"`
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

type EmailConfigMap map[string]string
type SMSConfigMap map[string]string
type PushConfigMap map[string]string
type ChatConfigMap map[string]string

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
		EmailProviders: make(map[string]registry.EmailConfig),
		SMSProviders:   make(map[string]registry.SMSConfig),
		PushProviders:  make(map[string]registry.PushConfig),
		ChatProviders:  make(map[string]registry.ChatConfig),
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	cfg.applyEnvOverrides()
	cfg.parseProviderConfigs()

	return cfg, nil
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
	buildCommon := func(m map[string]string) registry.CommonConfig {
		return registry.CommonConfig{
			APIKey:    m["api_key"],
			APISecret: m["api_secret"],
			Region:    m["region"],
			BaseURL:   m["base_url"],
			Extra:     m,
		}
	}

	for name, m := range c.Providers.Email {
		c.EmailProviders[name] = registry.EmailConfig{
			CommonConfig: buildCommon(map[string]string(m)),
			Domain:       m["domain"],
			FromEmail:    m["from_email"],
			FromName:     m["from_name"],
		}
	}

	for name, m := range c.Providers.SMS {
		c.SMSProviders[name] = registry.SMSConfig{
			CommonConfig: buildCommon(map[string]string(m)),
			FromPhone:    m["from_phone"],
		}
	}

	for name, m := range c.Providers.Push {
		c.PushProviders[name] = registry.PushConfig{
			CommonConfig: buildCommon(map[string]string(m)),
			AppID:        m["app_id"],
			Topic:        m["topic"],
		}
	}

	for name, m := range c.Providers.Chat {
		c.ChatProviders[name] = registry.ChatConfig{
			CommonConfig: buildCommon(map[string]string(m)),
			FromPhone:    m["from_phone"],
			WebhookURL:   m["webhook_url"],
		}
	}
}
