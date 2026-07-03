package contracts

// Attachment represents a file attachment for messages.
type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        []byte `json:"data,omitempty"`
	URL         string `json:"url,omitempty"`
}

const (
	MetaKeyProviderName = "provider_name"
	MetaKeyStoreContent = "store_content"

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

// StoreContentFromResult reports whether dispatch metadata requests persisting request/response bodies.
func StoreContentFromResult(r *SendResult) bool {
	if r == nil || r.Meta == nil {
		return false
	}
	return r.Meta[MetaKeyStoreContent] == metaStoreContentTrue
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
