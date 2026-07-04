package contracts

// Attachment represents a file attachment for messages.
type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        []byte `json:"data,omitempty"`
	URL         string `json:"url,omitempty"`
}

const (
	MetaKeyProviderName   = "provider_name"
	MetaKeyStoreContent   = "store_content"
	MetaKeyInboxMessageID = "inbox_message_id"
	MetaKeyDispatchMode   = "dispatch_mode"

	// Dispatch mode values mirror domain.MessageDispatchMode
	DispatchModeMemory   = "memory"
	DispatchModeProvider = "provider"

	metaStoreContentTrue  = "true"
	metaStoreContentFalse = "false"
)

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
	return r.Meta[MetaKeyProviderName]
}

// DispatchModeFromResult returns dispatch_mode from metadata when present.
func DispatchModeFromResult(r *SendResult) string {
	if r == nil || r.Meta == nil {
		return ""
	}
	return r.Meta[MetaKeyDispatchMode]
}

// StoreContentFromResult reports whether dispatch metadata requests persisting request/response bodies.
func StoreContentFromResult(r *SendResult) bool {
	if r == nil || r.Meta == nil {
		return false
	}
	return r.Meta[MetaKeyStoreContent] == metaStoreContentTrue
}

// InboxMessageIDFromResult returns the inbox message ID when dispatch captured content in the portal inbox.
//
// Provider dispatch stores the inbox UUID in MetaKeyInboxMessageID; SendResult.ID is the provider ID.
// Memory dispatch with store_message_content stores the inbox UUID in SendResult.ID.
func InboxMessageIDFromResult(r *SendResult) string {
	if r == nil || r.Meta == nil {
		return ""
	}
	if id := r.Meta[MetaKeyInboxMessageID]; id != "" {
		return id
	}
	if !StoreContentFromResult(r) {
		return ""
	}
	if r.Meta[MetaKeyDispatchMode] == DispatchModeMemory && r.ID != "" {
		return r.ID
	}
	return ""
}

// SetStoreContentMeta stamps store_content on r.Meta for downstream logging.
func SetStoreContentMeta(r *SendResult, storeContent bool) {
	if r.Meta == nil {
		r.Meta = make(map[string]string)
	}
	if storeContent {
		r.Meta[MetaKeyStoreContent] = metaStoreContentTrue
	} else {
		r.Meta[MetaKeyStoreContent] = metaStoreContentFalse
	}
}

// SetDispatchModeMeta stamps dispatch_mode on r.Meta.
func SetDispatchModeMeta(r *SendResult, mode string) {
	if r.Meta == nil {
		r.Meta = make(map[string]string)
	}
	r.Meta[MetaKeyDispatchMode] = mode
}
