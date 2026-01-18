package gateway

// Import internal/app to trigger provider registration via its imports.go.
// This ensures pkg/gateway uses the same providers as the server.
// To add new providers, update internal/app/imports.go — NOT this file.
import _ "github.com/weprodev/wpd-message-gateway/internal/app"
