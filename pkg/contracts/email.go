package contracts

import "context"

// Email represents an email message to be sent.
type Email struct {
	From        string            `json:"from,omitempty"`
	FromName    string            `json:"from_name,omitempty"`
	To          []string          `json:"to"`
	CC          []string          `json:"cc,omitempty"`
	BCC         []string          `json:"bcc,omitempty"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	PlainText   string            `json:"plain_text,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// EmailSender defines the contract for sending emails.
type EmailSender interface {
	Send(ctx context.Context, email *Email) (*SendResult, error)
	Name() string
}
