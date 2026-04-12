package domain

// Attachment represents a file attachment for messages.
type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        []byte `json:"data,omitempty"`
	URL         string `json:"url,omitempty"`
}

// SendResult represents the result of sending a message.
type SendResult struct {
	ID         string            `json:"id"`
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Meta       map[string]string `json:"meta,omitempty"`
}

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

// SMS represents an SMS message to be sent.
type SMS struct {
	From    string   `json:"from,omitempty"`
	To      []string `json:"to"`
	Message string   `json:"message"`
}

// PushNotification represents a push notification to be sent.
type PushNotification struct {
	DeviceTokens []string          `json:"device_tokens"`
	Title        string            `json:"title"`
	Body         string            `json:"body"`
	Data         map[string]string `json:"data,omitempty"`
	Badge        *int              `json:"badge,omitempty"`
	Sound        string            `json:"sound,omitempty"`
}

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
