// Package manager provides the central gateway for message dispatching.
//
// The Manager handles provider registration, initialization, and message routing
// for all message types (Email, SMS, Push, Chat).
//
// # Usage
//
//	cfg, _ := config.LoadFromEnv()
//	mgr, err := manager.New(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Send using default provider
//	result, err := mgr.SendEmail(ctx, &contracts.Email{...})
//
//	// Send using specific provider
//	result, err := mgr.SendEmailWith(ctx, "sendgrid", &contracts.Email{...})
//
// # Thread Safety
//
// Manager is safe for concurrent use. Provider maps are protected by sync.RWMutex.
//
// # Custom Providers
//
// Register custom provider implementations at runtime:
//
//	mgr.RegisterEmailProvider("custom", myCustomProvider)
package manager
