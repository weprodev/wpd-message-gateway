package postgres

import "fmt"

const (
	defaultListLimit = 50
	maxListLimit     = 500
)

// clampLimitOffset bounds list query pagination inputs.
func clampLimitOffset(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = defaultListLimit
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

// limitOffsetSQL returns a safe LIMIT/OFFSET clause for Postgres list queries.
// Values are clamped in Go — not passed as query parameters (driver limitation).
func limitOffsetSQL(limit, offset int) string {
	limit, offset = clampLimitOffset(limit, offset)
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}
