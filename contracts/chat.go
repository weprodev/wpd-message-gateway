package contracts

import "context"

// ChatMessage represents a message to be sent via chat/social platforms
// (WhatsApp, Telegram, Facebook Messenger, etc.)
type ChatMessage struct {
	From           string            `json:"from,omitempty"`            // Sender ID or phone number
	To             []string          `json:"to"`                        // Recipient IDs or phone numbers
	Message        string            `json:"message"`                   // Text message content
	Platform       string            `json:"platform,omitempty"`        // Platform: whatsapp, telegram, etc.
	TemplateID     string            `json:"template_id,omitempty"`     // Optional: template ID for WhatsApp Business API
	TemplateParams []string          `json:"template_params,omitempty"` // Optional: template parameters
	MediaURL       string            `json:"media_url,omitempty"`       // Optional: URL to media (image, video, document)
	MediaType      string            `json:"media_type,omitempty"`      // Optional: "image", "video", "audio", "document"
	Buttons        []ChatButton      `json:"buttons,omitempty"`         // Optional: interactive buttons
	ReplyToID      string            `json:"reply_to_id,omitempty"`     // Optional: message ID to reply to
	Metadata       map[string]string `json:"metadata,omitempty"`        // Optional: custom metadata
}

// ChatButton represents an interactive button in a chat message
type ChatButton struct {
	ID    string `json:"id"`              // Button identifier
	Text  string `json:"text"`            // Button display text
	URL   string `json:"url,omitempty"`   // Optional: URL for link buttons
	Phone string `json:"phone,omitempty"` // Optional: phone number for call buttons
}

// ChatSender defines the contract for sending chat/social media messages
type ChatSender interface {
	// Send sends a chat message and returns the result
	Send(ctx context.Context, message *ChatMessage) (*SendResult, error)

	// Name returns the provider name (e.g., "whatsapp", "telegram")
	Name() string
}
