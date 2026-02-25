package kavenegar

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/kavenegar/kavenegar-go"
	"github.com/nyaruka/phonenumbers"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const (
	ProviderName      = "kavenegar"
	defaultTimeout    = 40 * time.Second
	deduplicateWindow = 2 * time.Minute
)

type Config struct {
	APIKey    string
	FromPhone string
}

type Provider struct {
	client    *kavenegar.Kavenegar
	config    Config
	fromPhone string
	recent    map[string]time.Time
	mu        sync.Mutex
}

func New(cfg Config) (*Provider, error) {

	if cfg.APIKey == "" {
		return nil, errors.New("kavenegar: API Key is required")

	}

	if cfg.FromPhone == "" {
		return nil, errors.New("kavenegar: Sender number is required")

	}

	return &Provider{
		client:    kavenegar.New(cfg.APIKey),
		config:    cfg,
		fromPhone: cfg.FromPhone,
		recent:    make(map[string]time.Time),
	}, nil

}

func (p *Provider) Name() string {
	return ProviderName

}

func (p *Provider) Send(ctx context.Context, sms *contracts.SMS) (*contracts.SendResult, error) {
	if len(sms.To) == 0 {
		return nil, errors.New("at least one recipient is required")

	}

	messageLength := utf8.RuneCountInString(sms.Message)

	if messageLength > 160 {
		return nil, errors.New("message is too long, Limit is 160 characters")

	}

	for _, recipient := range sms.To {

		num, err := phonenumbers.Parse(recipient, "")

		if err != nil {
			return nil, fmt.Errorf("failed to parse phone number %s: %w", recipient, err)

		}

		if !phonenumbers.IsValidNumber(num) {
			return nil, fmt.Errorf("the phone number %s is not valid", recipient)

		}
	}

	p.mu.Lock()
	now := time.Now()

	for k, t := range p.recent {
		if now.Sub(t) > deduplicateWindow {
			delete(p.recent, k)
		}
	}

	for _, recipient := range sms.To {
		key := recipient + "|" + sms.Message

		if t, exists := p.recent[key]; exists {
			if now.Sub(t) < deduplicateWindow {
				p.mu.Unlock()
				return nil, fmt.Errorf("duplicate message detected for %s", recipient)
			}
		}

		p.recent[key] = now
	}
	p.mu.Unlock()

	sendCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	result, err := p.client.Message.Send(sms.From, sms.To, sms.Message, nil)

	if err != nil {
		if errors.Is(sendCtx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("request timed out after %v", defaultTimeout)
		}
		return nil, fmt.Errorf("kavenegar API error: %w", err)
	}

	response := &contracts.SendResult{
		ID:         fmt.Sprintf("%d", result[0].MessageID),
		StatusCode: 200,
		Message:    "SMS sent successfully",
	}
	return response, nil
}
