// Package contracts defines the interfaces and data types for message providers.
//
// This package contains:
//   - Message types (Email, SMS, PushNotification, ChatMessage)
//   - Sender interfaces (EmailSender, SMSSender, PushSender, ChatSender)
//   - Common types (Attachment, SendResult)
//
// # Interface Design
//
// All sender interfaces follow a consistent pattern:
//
//	type XxxSender interface {
//	    Send(ctx context.Context, message *Xxx) (*SendResult, error)
//	    Name() string
//	}
//
// This design allows:
//   - Easy provider switching
//   - Consistent error handling
//   - Context-based timeout/cancellation
//
// # Usage
//
//	func sendWelcomeEmail(sender contracts.EmailSender) error {
//	    _, err := sender.Send(ctx, &contracts.Email{
//	        To:      []string{"user@example.com"},
//	        Subject: "Welcome!",
//	        HTML:    "<h1>Hello!</h1>",
//	    })
//	    return err
//	}
package contracts
