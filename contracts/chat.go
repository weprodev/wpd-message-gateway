package contracts

import "context"

// ChatMessage represents a message to be sent via chat/social platforms
// (WhatsApp, Telegram, Facebook Messenger, etc.)
type ChatMessage struct {
	From           string            // Sender ID or phone number
	To             []string          // Recipient IDs or phone numbers
	Message        string            // Text message content
	TemplateID     string            // Optional: template ID for WhatsApp Business API
	TemplateParams []string          // Optional: template parameters
	MediaURL       string            // Optional: URL to media (image, video, document)
	MediaType      string            // Optional: "image", "video", "audio", "document"
	Buttons        []ChatButton      // Optional: interactive buttons
	ReplyToID      string            // Optional: message ID to reply to
	Metadata       map[string]string // Optional: custom metadata
}

// ChatButton represents an interactive button in a chat message
type ChatButton struct {
	ID    string // Button identifier
	Text  string // Button display text
	URL   string // Optional: URL for link buttons
	Phone string // Optional: phone number for call buttons
}

// ChatSender defines the contract for sending chat/social media messages
type ChatSender interface {
	// Send sends a chat message and returns the result
	Send(ctx context.Context, message *ChatMessage) (*SendResult, error)

	// Name returns the provider name (e.g., "whatsapp", "telegram")
	Name() string
}
