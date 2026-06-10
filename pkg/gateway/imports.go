// Package gateway imports all default providers to trigger their init() registration.
package gateway

import (
	// Register default providers
	_ "github.com/weprodev/wpd-message-gateway/pkg/provider/mailgun"
	_ "github.com/weprodev/wpd-message-gateway/pkg/provider/memory"
)
