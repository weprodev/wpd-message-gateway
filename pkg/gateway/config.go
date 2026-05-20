package gateway

import (
	"fmt"

	"github.com/weprodev/wpd-message-gateway/internal/registry"
)

// CommonConfig holds credential fields shared across all provider channels.
type CommonConfig struct {
	APIKey    string
	APISecret string
	Region    string
	BaseURL   string
}

// EmailConfig holds email-specific static provider configuration.
type EmailConfig struct {
	CommonConfig
	Domain    string
	FromEmail string
	FromName  string
}

// SMSConfig holds SMS-specific static provider configuration.
type SMSConfig struct {
	CommonConfig
	FromPhone string
}

// PushConfig holds push-notification-specific static provider configuration.
type PushConfig struct {
	CommonConfig
	AppID string
	Topic string
}

// ChatConfig holds chat-specific static provider configuration.
type ChatConfig struct {
	CommonConfig
	FromPhone  string
	WebhookURL string
}

// OTPConfig holds OTP-specific static provider configuration.
type OTPConfig struct {
	CommonConfig
	PhoneNumber    string
	Code           string
	CodeLength     string
	SenderUsername string
}

// Config is the static gateway configuration used with New().
// Providers are identified by name and must have been registered via init().
type Config struct {
	DefaultEmailProvider string
	DefaultSMSProvider   string
	DefaultPushProvider  string
	DefaultChatProvider  string
	DefaultOTPProvider   string

	EmailProviders map[string]EmailConfig
	SMSProviders   map[string]SMSConfig
	PushProviders  map[string]PushConfig
	ChatProviders  map[string]ChatConfig
	OTPProviders   map[string]OTPConfig
}

// buildEmailSenders constructs email senders from cfg and stores them in g.
func buildEmailSenders(g *Gateway, cfg Config) error {
	for name, ec := range cfg.EmailProviders {
		factory, err := registry.GetEmailFactory(name)
		if err != nil {
			return fmt.Errorf("gateway: email provider %q: %w", name, err)
		}
		sender, err := factory(toRegistryEmail(ec))
		if err != nil {
			return fmt.Errorf("gateway: init email provider %q: %w", name, err)
		}
		g.emailSenders[name] = sender
	}
	return nil
}

// buildSMSSenders constructs SMS senders from cfg and stores them in g.
func buildSMSSenders(g *Gateway, cfg Config) error {
	for name, sc := range cfg.SMSProviders {
		factory, err := registry.GetSMSFactory(name)
		if err != nil {
			return fmt.Errorf("gateway: SMS provider %q: %w", name, err)
		}
		sender, err := factory(toRegistrySMS(sc))
		if err != nil {
			return fmt.Errorf("gateway: init SMS provider %q: %w", name, err)
		}
		g.smsSenders[name] = sender
	}
	return nil
}

// buildPushSenders constructs push senders from cfg and stores them in g.
func buildPushSenders(g *Gateway, cfg Config) error {
	for name, pc := range cfg.PushProviders {
		factory, err := registry.GetPushFactory(name)
		if err != nil {
			return fmt.Errorf("gateway: push provider %q: %w", name, err)
		}
		sender, err := factory(toRegistryPush(pc))
		if err != nil {
			return fmt.Errorf("gateway: init push provider %q: %w", name, err)
		}
		g.pushSenders[name] = sender
	}
	return nil
}

// buildChatSenders constructs chat senders from cfg and stores them in g.
func buildChatSenders(g *Gateway, cfg Config) error {
	for name, cc := range cfg.ChatProviders {
		factory, err := registry.GetChatFactory(name)
		if err != nil {
			return fmt.Errorf("gateway: chat provider %q: %w", name, err)
		}
		sender, err := factory(toRegistryChat(cc))
		if err != nil {
			return fmt.Errorf("gateway: init chat provider %q: %w", name, err)
		}
		g.chatSenders[name] = sender
	}
	return nil
}

// buildOTPSenders constructs OTP senders from cfg and stores them in g.
func buildOTPSenders(g *Gateway, cfg Config) error {
	for name, oc := range cfg.OTPProviders {
		factory, err := registry.GetOTPFactory(name)
		if err != nil {
			return fmt.Errorf("gateway: OTP provider %q: %w", name, err)
		}
		sender, err := factory(toRegistryOTP(oc))
		if err != nil {
			return fmt.Errorf("gateway: init OTP provider %q: %w", name, err)
		}
		g.otpSenders[name] = sender
	}
	return nil
}

func toRegistryEmail(ec EmailConfig) registry.EmailConfig {
	return registry.EmailConfig{
		CommonConfig: registry.CommonConfig{APIKey: ec.APIKey, APISecret: ec.APISecret, Region: ec.Region, BaseURL: ec.BaseURL},
		Domain:       ec.Domain,
		FromEmail:    ec.FromEmail,
		FromName:     ec.FromName,
	}
}

func toRegistrySMS(sc SMSConfig) registry.SMSConfig {
	return registry.SMSConfig{
		CommonConfig: registry.CommonConfig{APIKey: sc.APIKey, APISecret: sc.APISecret, Region: sc.Region, BaseURL: sc.BaseURL},
		FromPhone:    sc.FromPhone,
	}
}

func toRegistryPush(pc PushConfig) registry.PushConfig {
	return registry.PushConfig{
		CommonConfig: registry.CommonConfig{APIKey: pc.APIKey, APISecret: pc.APISecret, Region: pc.Region, BaseURL: pc.BaseURL},
		AppID:        pc.AppID,
		Topic:        pc.Topic,
	}
}

func toRegistryChat(cc ChatConfig) registry.ChatConfig {
	return registry.ChatConfig{
		CommonConfig: registry.CommonConfig{APIKey: cc.APIKey, APISecret: cc.APISecret, Region: cc.Region, BaseURL: cc.BaseURL},
		FromPhone:    cc.FromPhone,
		WebhookURL:   cc.WebhookURL,
	}
}

func toRegistryOTP(oc OTPConfig) registry.OTPConfig {
	return registry.OTPConfig{
		CommonConfig:   registry.CommonConfig{APIKey: oc.APIKey, APISecret: oc.APISecret, Region: oc.Region, BaseURL: oc.BaseURL},
		PhoneNumber:    oc.PhoneNumber,
		Code:           oc.Code,
		CodeLength:     oc.CodeLength,
		SenderUsername: oc.SenderUsername,
	}
}
