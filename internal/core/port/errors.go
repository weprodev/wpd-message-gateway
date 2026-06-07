// Package port defines the interfaces (ports) that the core domain exposes
// to the outside world, and the shared error sentinels used across all layers.
package port

import "errors"

// Sentinel errors for repository operations.
// Use errors.Is(err, port.ErrNotFound) in callers — do not match on string.
var (
	// ErrNotFound is returned when a requested resource does not exist.
	ErrNotFound = errors.New("not found")

	// ErrConflict is returned when a create/update violates a uniqueness constraint.
	ErrConflict = errors.New("already exists")

	// ErrUnauthorized is returned when credentials are missing or invalid.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrInvalidCredentials is returned for failed portal login (unknown user or wrong password).
	// Callers may map it to HTTP 401 without treating it as an infrastructure failure.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrInvalidInput is returned when caller-supplied data fails validation.
	ErrInvalidInput = errors.New("invalid input")

	// ErrInternal is returned for unexpected infrastructure failures that are not
	// caused by the caller (e.g. DB connection drops, decryption failures).
	ErrInternal = errors.New("internal error")
)
