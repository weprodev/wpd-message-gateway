package kavenegar

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/kavenegar/kavenegar-go"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const (
	ProviderName = "kavenegar"
)

// Config holds Kavenegar-specific configuration.
type Config struct {
	APIKey    string
	FromPhone string
	BaseURL   string
}

// Provider implements port.SMSSender for Kavenegar.
type Provider struct {
	client    *kavenegar.Kavenegar
	config    Config
	fromPhone string
	dedup     *messageDeduplicator
}

type messageDeduplicator struct {
	mu   sync.Mutex
	sent map[string]struct{}
}

func newMessageDeduplicator() *messageDeduplicator {
	return &messageDeduplicator{
		sent: make(map[string]struct{}),
	}
}

func (d *messageDeduplicator) makeKey(recipient, message string) string {
	return recipient + "\x00" + message
}

func (d *messageDeduplicator) HasDuplicate(recipients []string, message string) (string, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, r := range recipients {
		key := d.makeKey(r, message)
		if _, exists := d.sent[key]; exists {
			return r, true
		}
	}

	return "", false
}

func (d *messageDeduplicator) MarkSent(recipients []string, message string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, r := range recipients {
		key := d.makeKey(r, message)
		d.sent[key] = struct{}{}
	}
}

var iranMobileRegexp = regexp.MustCompile(`^(?:\+98|0098|98|0)?9\d{9}$`)

func isValidMobile(phone string) bool {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return false
	}

	return iranMobileRegexp.MatchString(phone)
}

func New(cfg Config) (*Provider, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("kavenegar: API key is required")
	}

	api := kavenegar.New(cfg.APIKey)

	return &Provider{
		client:    api,
		config:    cfg,
		fromPhone: cfg.FromPhone,
		dedup:     newMessageDeduplicator(),
	}, nil
}

func (p *Provider) Name() string {
	return ProviderName
}

func (p *Provider) Send(ctx context.Context, sms *contracts.SMS) (*contracts.SendResult, error) {
	if len(sms.To) == 0 {
		return nil, errors.New("kavenegar: no recipients specified")
	}

	sender := sms.From
	usedDefaultSender := false

	if sender == "" {
		sender = p.fromPhone
		usedDefaultSender = true
	}

	if sender == "" {
		return nil, errors.New("kavenegar: no sender phone number specified and no default configured")
	}

	for _, recipient := range sms.To {
		if strings.TrimSpace(recipient) == "" {
			return nil, errors.New("kavenegar: recipient mobile number is required")
		}

		if !isValidMobile(recipient) {
			return nil, fmt.Errorf("kavenegar: recipient phone number %q is not a valid mobile number", recipient)
		}
	}

	if dupRecipient, ok := p.dedup.HasDuplicate(sms.To, sms.Message); ok {
		return nil, fmt.Errorf("kavenegar: duplicate message to recipient %s is not allowed", dupRecipient)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	res, err := p.client.Message.Send(sender, sms.To, sms.Message, nil)
	if err != nil {
		return nil, p.handleError(err)
	}

	if len(res) == 0 {
		return nil, errors.New("kavenegar: no result returned from API")
	}

	p.dedup.MarkSent(sms.To, sms.Message)

	firstResult := res[0]
	msg := "SMS sent successfully"

	meta := map[string]string{
		"message_id": fmt.Sprintf("%d", firstResult.MessageID),
		"status":     fmt.Sprintf("%d", firstResult.Status),
	}

	if usedDefaultSender {
		msg = fmt.Sprintf("%s (Used Default Sender: %s)", msg, sender)
		meta["default_sender_used"] = sender
	}

	return &contracts.SendResult{
		ID:         fmt.Sprintf("%d", firstResult.MessageID),
		StatusCode: 200,
		Message:    msg,
		Meta:       meta,
	}, nil
}

func (p *Provider) handleError(err error) error {
	switch e := err.(type) {
	case *kavenegar.APIError:
		return fmt.Errorf("kavenegar:\n%w", e)
	case *kavenegar.HTTPError:
		return fmt.Errorf("kavenegar:\nHTTP error: %w", e)
	default:
		return fmt.Errorf("kavenegar:\n%w", e)
	}
}
