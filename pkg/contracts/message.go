package contracts

// Attachment represents a file attachment for messages.
type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        []byte `json:"data,omitempty"`
	URL         string `json:"url,omitempty"`
}

// SendResultItem holds per-recipient result for batch sends.
type SendResultItem struct {
	PhoneNumber string `json:"phone_number"`
	RequestID   string `json:"request_id,omitempty"`
	StatusCode  int    `json:"status_code"`
	Error       string `json:"error"`
}

// SendResult represents the result of sending a message.
type SendResult struct {
	ID         string            `json:"id"`
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Meta       map[string]string `json:"meta,omitempty"`
	Items      []SendResultItem  `json:"items,omitempty"`
}
