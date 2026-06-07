# WPD Message Gateway SDK

The `pkg` directory contains the core SDK for the WPD Message Gateway. It is fully independent from the portal and backend infrastructure, providing a clean, decoupled interface for 3rd party clients and external systems to integrate with multiple messaging providers (Email, SMS, Push, Chat).

## Architecture & Concepts

The SDK adheres to clean architecture, Domain-Driven Design (DDD), and SOLID principles:

- **`contracts/`**: Contains the pure domain interfaces (`EmailSender`, `SMSSender`, etc.) and data models (`Email`, `SMS`, `SendResult`). This acts as the universal language across all providers.
- **`provider/`**: Contains the concrete implementations of those contracts for specific services (e.g., `mailgun`). It also includes a `memory` provider for mocking, testing, and gateway interception.
- **`registry/`**: A central factory registry that allows dynamic registration and instantiation of providers based on arbitrary configuration maps.
- **`gateway/`**: Exposes generic client definitions for REST integrations if needed.

## How to Use as a Client

Any external system or 3rd party application can import this package to seamlessly send messages through any supported provider without worrying about the underlying HTTP clients or vendor SDKs.

### 1. Direct Provider Usage

If you know exactly which provider you want to use, you can instantiate it directly using the contract interface:

```go
package main

import (
	"context"
	"log"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
	"github.com/weprodev/wpd-message-gateway/pkg/provider/mailgun"
)

func main() {
	// 1. Configure the specific provider
	sender, err := mailgun.New(mailgun.Config{
		Domain:    "your-domain.com",
		APIKey:    "your-api-key",
		FromEmail: "no-reply@your-domain.com",
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Build the standard contract message
	email := &contracts.Email{
		To:      []string{"user@example.com"},
		Subject: "Welcome to WPD",
		HTML:    "<h1>Hello World</h1>",
		From:    "no-reply@your-domain.com",
	}

	// 3. Send using the unified interface
	result, err := sender.Send(context.Background(), email)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("Sent! ID: %s", result.ID)
}
```

### 2. Dynamic Provider Resolution (Registry)

For systems that store configurations dynamically (e.g., in a database), you can use the Registry. This allows you to swap providers at runtime without changing your code:

```go
package main

import (
	"context"
	"log"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
	"github.com/weprodev/wpd-message-gateway/pkg/registry"
	_ "github.com/weprodev/wpd-message-gateway/pkg/provider/mailgun" // Import for side-effects (registration)
)

func main() {
	// Configuration loaded from your DB or Environment
	config := map[string]interface{}{
		"domain":  "your-domain.com",
		"api_key": "your-api-key",
	}

	// Instantiate the provider dynamically by name
	sender, err := registry.BuildEmailSender("mailgun", config)
	if err != nil {
		log.Fatal(err)
	}

	email := &contracts.Email{
		To:      []string{"user@example.com"},
		Subject: "Dynamic Send",
		HTML:    "<p>Works perfectly.</p>",
		From:    "no-reply@your-domain.com",
	}

	result, err := sender.Send(context.Background(), email)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Creating a New Provider

To add a new provider (e.g., Twilio for SMS):
1. Create a new folder in `pkg/provider/twilio`.
2. Implement the relevant contract (`contracts.SMSSender`).
3. Add a builder function and register it via `registry.RegisterSMSBuilder("twilio", builderFn)` in an `init()` block.
