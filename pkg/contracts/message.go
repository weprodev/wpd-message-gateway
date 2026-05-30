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

// DeliveryStatus represents the current message delivery status.
type DeliveryStatus struct {
	Status        string `json:"status"`
	LastUpdatedAt int64  `json:"last_updated_at,omitempty"`
}

// VerificationProcessStatus represents the current status of the verification process.
type VerificationProcessStatus struct {
	Status     string `json:"status"`
	VerifiedAt int64  `json:"verified_at,omitempty"`
}

// VerificationStatus represents the response from checkVerificationStatus.
type VerificationStatus struct {
	RequestID          string                     `json:"request_id"`
	PhoneNumber        string                     `json:"phone_number,omitempty"`
	RequestCost        float64                    `json:"request_cost"`
	IsRefunded         *bool                      `json:"is_refunded,omitempty"`
	RemainingBalance   *float64                   `json:"remaining_balance,omitempty"`
	DeliveryStatus     *DeliveryStatus            `json:"delivery_status,omitempty"`
	VerificationStatus *VerificationProcessStatus `json:"verification_status,omitempty"`
}

// SendResult represents the result of sending a message.
type SendResult struct {
	ID         string            `json:"id"`
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Meta       map[string]string `json:"meta,omitempty"`
	Items      []SendResultItem  `json:"items,omitempty"`
}
