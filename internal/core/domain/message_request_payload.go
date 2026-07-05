package domain

import "time"

// MessageRequestPayload stores the body of a message request separately
// from the main metadata logs for optimization and privacy control.
type MessageRequestPayload struct {
	LogID        string
	RequestBody  string
	ResponseBody string
	CreatedAt    time.Time
}
