package contracts

import "context"

// Email represents an email message to be sent
type Email struct {
	From        string            // Sender email address
	FromName    string            // Sender display name
	To          []string          // Recipients
	CC          []string          // Carbon copy recipients
	BCC         []string          // Blind carbon copy recipients
	ReplyTo     string            // Reply-to address
	Subject     string            // Email subject
	HTML        string            // HTML body content
	PlainText   string            // Plain text body content
	Attachments []Attachment      // File attachments
	Headers     map[string]string // Custom email headers
}

// EmailSender defines the contract for sending emails
type EmailSender interface {
	// Send sends an email and returns the result
	Send(ctx context.Context, email *Email) (*SendResult, error)

	// Name returns the provider name (e.g., "mailgun", "sendgrid")
	Name() string
}
