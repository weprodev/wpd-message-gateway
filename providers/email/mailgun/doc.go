// Package mailgun implements the [contracts.EmailSender] interface using the Mailgun API.
//
// # Configuration
//
// Configure this provider using the following environment variables:
//
//   - MESSAGE_MAILGUN_API_KEY    (Required)
//   - MESSAGE_MAILGUN_DOMAIN     (Required)
//   - MESSAGE_MAILGUN_FROM_EMAIL (Optional, default sender)
//   - MESSAGE_MAILGUN_BASE_URL   (Optional, e.g. "https://api.eu.mailgun.net")
//
// # Usage
//
// This package is typically used via the Manager, not directly.
//
//	// The manager will automatically initialize this provider if configured.
//	mgr.SendEmail(ctx, email)
package mailgun
