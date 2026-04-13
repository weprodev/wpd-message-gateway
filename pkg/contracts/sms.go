package contracts

// SMS represents an SMS message to be sent.
type SMS struct {
	From    string   `json:"from,omitempty"`
	To      []string `json:"to"`
	Message string   `json:"message"`
}
