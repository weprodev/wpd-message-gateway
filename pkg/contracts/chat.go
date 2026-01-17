package contracts

import "context"

// ChatMessage represents a message to be sent via chat/social platforms.
type ChatMessage struct {
	From           string            `json:"from,omitempty"`
	To             []string          `json:"to"`
	Message        string            `json:"message"`
	Platform       string            `json:"platform,omitempty"`
	TemplateID     string            `json:"template_id,omitempty"`
	TemplateParams []string          `json:"template_params,omitempty"`
	MediaURL       string            `json:"media_url,omitempty"`
	MediaType      string            `json:"media_type,omitempty"`
	Buttons        []ChatButton      `json:"buttons,omitempty"`
	ReplyToID      string            `json:"reply_to_id,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// ChatButton represents an interactive button in a chat message.
type ChatButton struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	URL   string `json:"url,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// ChatSender defines the contract for sending chat/social media messages.
type ChatSender interface {
	Send(ctx context.Context, message *ChatMessage) (*SendResult, error)
	Name() string
}
