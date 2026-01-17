// Package errors provides structured error types for the message gateway.
package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors for common failure cases.
var (
	ErrNoRecipients      = errors.New("no recipients specified")
	ErrNoFromAddress     = errors.New("no from address specified")
	ErrNoDefaultProvider = errors.New("no default provider configured")
)

// ProviderError represents an error from a message provider.
type ProviderError struct {
	Provider   string
	StatusCode int
	Message    string
	Err        error
}

func (e *ProviderError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (code: %d): %v", e.Provider, e.Message, e.StatusCode, e.Err)
	}
	return fmt.Sprintf("%s: %s (code: %d)", e.Provider, e.Message, e.StatusCode)
}

func (e *ProviderError) Unwrap() error {
	return e.Err
}

// NewProviderError creates a new ProviderError.
func NewProviderError(provider, message string, statusCode int, err error) *ProviderError {
	return &ProviderError{
		Provider:   provider,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// ConfigError represents a configuration error.
type ConfigError struct {
	Provider string
	Field    string
	Message  string
	Err      error
}

func (e *ConfigError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("config error for %s: %s - %s: %v", e.Provider, e.Field, e.Message, e.Err)
	}
	return fmt.Sprintf("config error for %s: %s - %s", e.Provider, e.Field, e.Message)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// NewConfigError creates a new ConfigError.
func NewConfigError(provider, field, message string) *ConfigError {
	return &ConfigError{
		Provider: provider,
		Field:    field,
		Message:  message,
	}
}

// ProviderNotFoundError indicates a requested provider doesn't exist.
type ProviderNotFoundError struct {
	ProviderType string
	ProviderName string
}

func (e *ProviderNotFoundError) Error() string {
	return fmt.Sprintf("%s provider '%s' not found", e.ProviderType, e.ProviderName)
}

// NewProviderNotFoundError creates a new ProviderNotFoundError.
func NewProviderNotFoundError(providerType, providerName string) *ProviderNotFoundError {
	return &ProviderNotFoundError{
		ProviderType: providerType,
		ProviderName: providerName,
	}
}

// IsProviderNotFound checks if an error is a ProviderNotFoundError.
func IsProviderNotFound(err error) bool {
	var pnf *ProviderNotFoundError
	return errors.As(err, &pnf)
}

// IsConfigError checks if an error is a ConfigError.
func IsConfigError(err error) bool {
	var ce *ConfigError
	return errors.As(err, &ce)
}

// IsProviderError checks if an error is a ProviderError.
func IsProviderError(err error) bool {
	var pe *ProviderError
	return errors.As(err, &pe)
}
