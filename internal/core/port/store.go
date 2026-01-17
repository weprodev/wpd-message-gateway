package port

// MessageStore is the interface for message storage (used by DevBox).
// This allows providers to store messages without circular imports.
type MessageStore interface {
	// Marker interface - concrete methods are on the implementation
	// Providers that need storage cast this to their specific type
}
