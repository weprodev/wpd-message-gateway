package contracts

// Attachment represents a file attachment for messages
type Attachment struct {
	Filename    string // Name of the file
	ContentType string // MIME type of the file
	Data        []byte // File content
}

// SendResult represents the result of sending a message
type SendResult struct {
	ID         string            // Provider-specific message ID
	StatusCode int               // HTTP status code from the provider
	Message    string            // Human-readable status message
	Meta       map[string]string // Additional provider-specific metadata
}
