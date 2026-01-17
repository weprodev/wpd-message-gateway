package mailgun

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/contracts"
	msgerrors "github.com/weprodev/wpd-message-gateway/errors"
)

const (
	ProviderName   = "mailgun"
	defaultTimeout = 30 * time.Second
)

// Provider implements contracts.EmailSender for Mailgun
type Provider struct {
	client      *mailgun.MailgunImpl
	config      config.EmailConfig
	fromAddress string
	fromName    string
}

// New creates a new Mailgun provider from configuration
func New(cfg config.EmailConfig) (*Provider, error) {
	if cfg.APIKey == "" {
		return nil, msgerrors.NewConfigError(ProviderName, "APIKey", "API key is required")
	}
	if cfg.Domain == "" {
		return nil, msgerrors.NewConfigError(ProviderName, "Domain", "domain is required")
	}

	mg := mailgun.NewMailgun(cfg.Domain, cfg.APIKey)

	// Set EU endpoint if configured
	if cfg.BaseURL != "" {
		baseURL := cfg.BaseURL
		// Mailgun requires the base URL to end with /v3 (or v1/v2/v4)
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

// Name returns the provider name
func (p *Provider) Name() string {
	return ProviderName
}

// Send sends an email via Mailgun
func (p *Provider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
	if len(email.To) == 0 {
		return nil, msgerrors.ErrNoRecipients
	}

	// Build sender address
	from := p.buildFromAddress(email)
	if from == "" {
		return nil, msgerrors.ErrNoFromAddress
	}

	// Create base message
	msg := mailgun.NewMessage(from, email.Subject, email.PlainText, email.To...)

	// Set HTML body if provided
	if email.HTML != "" {
		msg.SetHTML(email.HTML)
	}

	// Add CC recipients
	for _, cc := range email.CC {
		msg.AddCC(cc)
	}

	// Add BCC recipients
	for _, bcc := range email.BCC {
		msg.AddBCC(bcc)
	}

	// Set reply-to
	if email.ReplyTo != "" {
		msg.SetReplyTo(email.ReplyTo)
	}

	// Add custom headers
	for key, value := range email.Headers {
		msg.AddHeader(key, value)
	}

	// Add attachments
	for _, att := range email.Attachments {
		msg.AddBufferAttachment(att.Filename, att.Data)
	}

	// Send with timeout
	sendCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, id, err := p.client.Send(sendCtx, msg)
	if err != nil {
		return nil, msgerrors.NewProviderError(ProviderName, "failed to send email", 500, err)
	}

	return &contracts.SendResult{
		ID:         id,
		StatusCode: 200,
		Message:    "Email sent successfully",
	}, nil
}

// buildFromAddress constructs the from address string
func (p *Provider) buildFromAddress(email *contracts.Email) string {
	fromAddr := email.From
	fromName := email.FromName

	// Use default values if not specified
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
