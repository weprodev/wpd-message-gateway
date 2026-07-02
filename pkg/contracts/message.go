package contracts

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

// ProviderNameFromResult returns provider_name from dispatch metadata when present.
func ProviderNameFromResult(r *SendResult) string {
	if r == nil || r.Meta == nil {
		return ""
	}
	return r.Meta["provider_name"]
}
