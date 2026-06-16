package domain

import "time"

// StoredMessageDispatchStatus tracks provider dispatch outcome for a stored payload.
type StoredMessageDispatchStatus string

const (
	StoredMessageDispatchPending StoredMessageDispatchStatus = "pending"
	StoredMessageDispatchSent    StoredMessageDispatchStatus = "sent"
	StoredMessageDispatchFailed  StoredMessageDispatchStatus = "failed"
)

// StoredMessageDispatchOutcome captures the minimal provider result linked to a stored message.
type StoredMessageDispatchOutcome struct {
	Status             StoredMessageDispatchStatus
	ProviderMessageID  string
	ProviderStatusCode int // zero means unset
	DispatchError      string
	DispatchedAt       time.Time
}
