package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	pkgconfig "github.com/weprodev/go-pkg/config"

	"github.com/weprodev/wpd-message-gateway/pkg/registry"
)

// Config represents the application configuration.
type Config struct {
	Environment   string                   `yaml:"environment"`
	EncryptionKey string                   `yaml:"encryption_key,omitempty"`
	Server        ServerConfig             `yaml:"server"`
	Portal        PortalConfig             `yaml:"portal,omitempty"`
	Database      pkgconfig.DatabaseConfig `yaml:"database,omitempty"`
	Providers     ProviderConfig           `yaml:"providers"`
	// Parsed provider configs - using registry types as single source of truth
	EmailProviders map[string]registry.EmailConfig `yaml:"-"`
	SMSProviders   map[string]registry.SMSConfig   `yaml:"-"`
	PushProviders  map[string]registry.PushConfig  `yaml:"-"`
	ChatProviders  map[string]registry.ChatConfig  `yaml:"-"`
}

// ServerConfig holds server configuration.
type ServerConfig struct {
	Port int `yaml:"port"`
}

// PortalConfig holds JWT settings and portal URL for the web portal.
type PortalConfig struct {
	JWTSecret   string `yaml:"jwt_secret"`
	JWTTTLHours int    `yaml:"jwt_ttl_hours"`
	UIPort      int    `yaml:"ui_port"`
	BaseURL     string `yaml:"base_url"`
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

const defaultServerPort = 10101

// EnvironmentName returns the runtime environment (config file, then APP_ENVIRONMENT, then local).
func (c *Config) EnvironmentName() string {
	if c != nil && c.Environment != "" {
		return c.Environment
	}
	if v := os.Getenv("APP_ENVIRONMENT"); v != "" {
		return v
	}
	return "local"
}

// ListenAddr returns the HTTP listen address for the gateway server.
func (c *Config) ListenAddr() string {
	port := defaultServerPort
	if c != nil && c.Server.Port > 0 {
		port = c.Server.Port
	}
	return fmt.Sprintf(":%d", port)
}

// LogSummary logs a non-secret snapshot of loaded configuration at startup.
func (c *Config) LogSummary() {
	if c == nil {
		return
	}
	slog.Info("loaded configuration",
		"environment", c.EnvironmentName(),
		"listen_addr", c.ListenAddr(),
		"email_provider", c.DefaultEmailProvider(),
		"sms_provider", c.DefaultSMSProvider(),
		"push_provider", c.DefaultPushProvider(),
		"chat_provider", c.DefaultChatProvider(),
	)
}

// LoadAppConfig loads configuration from YAML, applying env overrides.
// When configPath is empty, resolves via CONFIG_PATH, container paths, or configs/local.yml.
func LoadAppConfig(configPath string) (*Config, error) {
	configPath = resolveConfigPath(configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", configPath, err)
	}

	cfg := &Config{
		EmailProviders: make(map[string]registry.EmailConfig),
		SMSProviders:   make(map[string]registry.SMSConfig),
		PushProviders:  make(map[string]registry.PushConfig),
		ChatProviders:  make(map[string]registry.ChatConfig),
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file %s: %w", configPath, err)
	}

	cfg.applyEnvOverrides()
	cfg.parseProviderConfigs()

	return cfg, nil
}

// resolveConfigPath returns an explicit path or discovers one from env/container defaults.
func resolveConfigPath(configPath string) string {
	if configPath != "" {
		return configPath
	}
	if v := os.Getenv("CONFIG_PATH"); v != "" {
		return v
	}

	env := os.Getenv("APP_ENVIRONMENT")
	if env == "" {
		env = os.Getenv("ENVIRONMENT")
	}
	if env == "" {
		env = "local"
	}

	for _, candidate := range []string{
		fmt.Sprintf("/configs/%s.yml", env),
		"/configs/local.yml",
		"/configs/config.yml",
		"configs/local.yml",
	} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return "configs/local.yml"
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
		case "MESSAGE_JWT_SECRET":
			c.Portal.JWTSecret = val
		case "MESSAGE_PORTAL_BASE_URL":
			c.Portal.BaseURL = val
		case "MESSAGE_CONFIG_ENCRYPTION_KEY":
			c.EncryptionKey = val
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
