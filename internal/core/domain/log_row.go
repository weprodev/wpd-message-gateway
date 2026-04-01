package domain

// MessageRequestLogWithSource is a log row joined with API key metadata for portal lists.
type MessageRequestLogWithSource struct {
	MessageRequestLog
	SourceName string `json:"source_name"`
	ClientID   string `json:"client_id,omitempty"`
}
