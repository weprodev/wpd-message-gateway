package postgres

import (
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

func mapDBError(err error) error {
	if err == nil {
		return nil
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return fmt.Errorf("%w: %s", port.ErrConflict, pqErr.Constraint)
	}
	return err
}
