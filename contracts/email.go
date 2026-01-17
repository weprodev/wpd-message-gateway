package contracts

import "context"

// Email represents an email message to be sent
type Email struct {
	From        string            `json:"from,omitempty"`        // Sender email address
	FromName    string            `json:"from_name,omitempty"`   // Sender display name
	To          []string          `json:"to"`                    // Recipients
	CC          []string          `json:"cc,omitempty"`          // Carbon copy recipients
	BCC         []string          `json:"bcc,omitempty"`         // Blind carbon copy recipients
	ReplyTo     string            `json:"reply_to,omitempty"`    // Reply-to address
	Subject     string            `json:"subject"`               // Email subject
	HTML        string            `json:"html,omitempty"`        // HTML body content
	PlainText   string            `json:"plain_text,omitempty"`  // Plain text body content
	Attachments []Attachment      `json:"attachments,omitempty"` // File attachments
	Headers     map[string]string `json:"headers,omitempty"`     // Custom email headers
}

// EmailSender defines the contract for sending emails
type EmailSender interface {
	// Send sends an email and returns the result
	Send(ctx context.Context, email *Email) (*SendResult, error)

	// Name returns the provider name (e.g., "mailgun", "sendgrid")
	Name() string
}
