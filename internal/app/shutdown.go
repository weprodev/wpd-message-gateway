package app

import (
	"context"
	"fmt"
)

// Shutdown stops the HTTP server and closes the database connection pool.
func (a *Application) Shutdown(ctx context.Context) error {
	if a == nil {
		return nil
	}

	var shutdownErr error
	if a.Echo != nil {
		if err := a.Echo.Shutdown(ctx); err != nil {
			shutdownErr = fmt.Errorf("echo shutdown: %w", err)
		}
	}
	if a.PgClient != nil {
		if err := a.PgClient.Close(); err != nil {
			if shutdownErr != nil {
				return fmt.Errorf("%v; pg close: %w", shutdownErr, err)
			}
			return fmt.Errorf("pg close: %w", err)
		}
	}
	return shutdownErr
}
