package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors for common failure cases
var (
	// ErrNoRecipients indicates no recipients were specified
	ErrNoRecipients = errors.New("no recipients specified")

	// ErrNoFromAddress indicates no sender address was specified
	ErrNoFromAddress = errors.New("no from address specified")

	// ErrNoDefaultProvider indicates no default provider is configured
	ErrNoDefaultProvider = errors.New("no default provider configured")
)

// ProviderError represents an error from a message provider
type ProviderError struct {
	Provider   string // Provider name (e.g., "mailgun")
	StatusCode int    // HTTP status code if applicable
	Message    string // Human-readable error message
	Err        error  // Underlying error
}

// Error implements the error interface
func (e *ProviderError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (code: %d): %v", e.Provider, e.Message, e.StatusCode, e.Err)
	}
	return fmt.Sprintf("%s: %s (code: %d)", e.Provider, e.Message, e.StatusCode)
}

// Unwrap returns the underlying error for errors.Is/As support
func (e *ProviderError) Unwrap() error {
	return e.Err
}

// NewProviderError creates a new ProviderError
func NewProviderError(provider, message string, statusCode int, err error) *ProviderError {
	return &ProviderError{
		Provider:   provider,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// ConfigError represents a configuration error
type ConfigError struct {
	Provider string // Provider name
	Field    string // Configuration field with the issue
	Message  string // Error description
	Err      error  // Underlying error (optional)
}

// Error implements the error interface
func (e *ConfigError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("config error for %s: %s - %s: %v", e.Provider, e.Field, e.Message, e.Err)
	}
	return fmt.Sprintf("config error for %s: %s - %s", e.Provider, e.Field, e.Message)
}

// Unwrap returns the underlying error for errors.Is/As support
func (e *ConfigError) Unwrap() error {
	return e.Err
}

// NewConfigError creates a new ConfigError
func NewConfigError(provider, field, message string) *ConfigError {
	return &ConfigError{
		Provider: provider,
		Field:    field,
		Message:  message,
	}
}

// NewConfigErrorWithCause creates a ConfigError wrapping an underlying error
func NewConfigErrorWithCause(provider, field, message string, err error) *ConfigError {
	return &ConfigError{
		Provider: provider,
		Field:    field,
		Message:  message,
		Err:      err,
	}
}

// ProviderNotFoundError indicates a requested provider doesn't exist
type ProviderNotFoundError struct {
	ProviderType string // Type of provider (email, sms, push)
	ProviderName string // Name of the requested provider
}

// Error implements the error interface
func (e *ProviderNotFoundError) Error() string {
	return fmt.Sprintf("%s provider '%s' not found", e.ProviderType, e.ProviderName)
}

// NewProviderNotFoundError creates a new ProviderNotFoundError
func NewProviderNotFoundError(providerType, providerName string) *ProviderNotFoundError {
	return &ProviderNotFoundError{
		ProviderType: providerType,
		ProviderName: providerName,
	}
}
