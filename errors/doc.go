// Package errors provides error types for the message gateway.
//
// This package defines structured error types that provide:
//   - Provider identification
//   - HTTP status codes (when applicable)
//   - Error wrapping for errors.Is/As support
//
// # Error Types
//
//   - ProviderError: Errors from message providers (API failures, rate limits)
//   - ConfigError: Configuration validation errors
//   - ProviderNotFoundError: Requested provider doesn't exist
//
// # Usage
//
//	if err != nil {
//	    var provErr *errors.ProviderError
//	    if errors.As(err, &provErr) {
//	        log.Printf("Provider %s failed: %s", provErr.Provider, provErr.Message)
//	    }
//	}
package errors
