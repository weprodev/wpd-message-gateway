package contracts

// Attachment represents a file attachment for messages
type Attachment struct {
	Filename    string `json:"filename"`       // Name of the file
	ContentType string `json:"content_type"`   // MIME type of the file
	Data        []byte `json:"data,omitempty"` // File content
	URL         string `json:"url,omitempty"`  // URL to attachment (alternative to Data)
}

// SendResult represents the result of sending a message
type SendResult struct {
	ID         string            `json:"id"`             // Provider-specific message ID
	StatusCode int               `json:"status_code"`    // HTTP status code from the provider
	Message    string            `json:"message"`        // Human-readable status message
	Meta       map[string]string `json:"meta,omitempty"` // Additional provider-specific metadata
}
