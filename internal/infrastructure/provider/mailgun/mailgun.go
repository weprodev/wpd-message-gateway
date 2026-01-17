package mailgun

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const (
	ProviderName   = "mailgun"
	defaultTimeout = 30 * time.Second
)

// Config holds Mailgun-specific configuration.
type Config struct {
	APIKey    string
	Domain    string
	BaseURL   string
	FromEmail string
	FromName  string
}

// Provider implements port.EmailSender for Mailgun.
type Provider struct {
	client      *mailgun.MailgunImpl
	config      Config
	fromAddress string
	fromName    string
}

// New creates a new Mailgun provider.
func New(cfg Config) (*Provider, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("mailgun: API key is required")
	}
	if cfg.Domain == "" {
		return nil, errors.New("mailgun: domain is required")
	}

	mg := mailgun.NewMailgun(cfg.Domain, cfg.APIKey)

	if cfg.BaseURL != "" {
		baseURL := cfg.BaseURL
		if !strings.HasSuffix(baseURL, "/v3") &&
			!strings.HasSuffix(baseURL, "/v4") &&
			!strings.HasSuffix(baseURL, "/v1") {
			baseURL = strings.TrimRight(baseURL, "/") + "/v3"
		}
		mg.SetAPIBase(baseURL)
	}

	return &Provider{
		client:      mg,
		config:      cfg,
		fromAddress: cfg.FromEmail,
		fromName:    cfg.FromName,
	}, nil
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return ProviderName
}

// Send sends an email via Mailgun.
func (p *Provider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
	if len(email.To) == 0 {
		return nil, errors.New("no recipients specified")
	}

	from := p.buildFromAddress(email)
	if from == "" {
		return nil, errors.New("no from address specified")
	}

	msg := mailgun.NewMessage(from, email.Subject, email.PlainText, email.To...)

	if email.HTML != "" {
		msg.SetHTML(email.HTML)
	}

	for _, cc := range email.CC {
		msg.AddCC(cc)
	}

	for _, bcc := range email.BCC {
		msg.AddBCC(bcc)
	}

	if email.ReplyTo != "" {
		msg.SetReplyTo(email.ReplyTo)
	}

	for key, value := range email.Headers {
		msg.AddHeader(key, value)
	}

	for _, att := range email.Attachments {
		msg.AddBufferAttachment(att.Filename, att.Data)
	}

	sendCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, id, err := p.client.Send(sendCtx, msg)
	if err != nil {
		return nil, fmt.Errorf("mailgun: failed to send email: %w", err)
	}

	return &contracts.SendResult{
		ID:         id,
		StatusCode: 200,
		Message:    "Email sent successfully",
	}, nil
}

func (p *Provider) buildFromAddress(email *contracts.Email) string {
	fromAddr := email.From
	fromName := email.FromName

	if fromAddr == "" {
		fromAddr = p.fromAddress
	}
	if fromName == "" {
		fromName = p.fromName
	}

	if fromAddr == "" {
		return ""
	}

	if fromName != "" {
		return fmt.Sprintf("%s <%s>", fromName, fromAddr)
	}
	return fromAddr
}
