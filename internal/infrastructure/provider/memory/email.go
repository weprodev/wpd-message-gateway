package memory

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// MailpitConfig holds configuration for SMTP forwarding to Mailpit.
type MailpitConfig struct {
	Enabled bool
}

// --- SMTP Forwarder ---

const (
	defaultMailpitHost = "localhost"
	defaultMailpitPort = "10102"
)

type smtpForwarder struct {
	host    string
	port    string
	enabled bool
}

func newSMTPForwarder(cfg MailpitConfig) *smtpForwarder {
	return &smtpForwarder{
		host:    defaultMailpitHost,
		port:    defaultMailpitPort,
		enabled: cfg.Enabled,
	}
}

func (f *smtpForwarder) forward(email *contracts.Email) {
	if !f.enabled {
		return
	}

	from := email.From
	if from == "" {
		from = "devbox@local.dev"
	}

	if len(email.To) == 0 {
		return
	}

	msg := f.buildMessage(email, from)
	recipients := f.collectRecipients(email)

	addr := fmt.Sprintf("%s:%s", f.host, f.port)
	if err := smtp.SendMail(addr, nil, from, recipients, []byte(msg)); err != nil {
		log.Printf("SMTP forward to Mailpit failed: %v", err)
		return
	}

	log.Printf("Email forwarded to Mailpit: %s -> %v", email.Subject, email.To)
}

func (f *smtpForwarder) buildMessage(email *contracts.Email, from string) string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))

	if len(email.CC) > 0 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(email.CC, ", ")))
	}

	if email.HTML != "" {
		msg.WriteString("MIME-Version: 1.0\r\n")
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		msg.WriteString("\r\n")
		msg.WriteString(email.HTML)
	} else {
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		msg.WriteString("\r\n")
		msg.WriteString(email.PlainText)
	}

	return msg.String()
}

func (f *smtpForwarder) collectRecipients(email *contracts.Email) []string {
	recipients := make([]string, 0, len(email.To)+len(email.CC)+len(email.BCC))
	recipients = append(recipients, email.To...)
	recipients = append(recipients, email.CC...)
	recipients = append(recipients, email.BCC...)
	return recipients
}

// --- Email Provider ---

// EmailProvider implements port.EmailSender using an in-memory store.
type EmailProvider struct {
	store         *Store
	smtpForwarder *smtpForwarder
}

// NewEmailProvider creates a new memory email provider.
func NewEmailProvider(store *Store, mailpitCfg MailpitConfig) *EmailProvider {
	return &EmailProvider{
		store:         store,
		smtpForwarder: newSMTPForwarder(mailpitCfg),
	}
}

// Store returns the underlying memory store.
func (e *EmailProvider) Store() *Store {
	return e.store
}

// Name returns the provider name.
func (e *EmailProvider) Name() string {
	return ProviderName
}

// Send stores the email in memory and optionally forwards to Mailpit.
func (e *EmailProvider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
	id := uuid.New().String()

	stored := &StoredEmail{
		ID:        id,
		CreatedAt: time.Now(),
		Email:     email,
	}
	e.store.AddEmail(stored)

	if e.smtpForwarder != nil {
		e.smtpForwarder.forward(email)
	}

	return &contracts.SendResult{
		ID:         id,
		StatusCode: 200,
		Message:    "Stored email in memory",
	}, nil
}
