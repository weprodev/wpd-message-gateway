// Package app imports all providers to trigger their init() registration.
// Add new provider imports here — this is the ONLY file to modify when adding providers.
package app

import (
	// Built-in providers
	_ "github.com/weprodev/wpd-message-gateway/pkg/provider/mailgun"
	_ "github.com/weprodev/wpd-message-gateway/pkg/provider/memory"
)
