package postgres

import "database/sql"

func optionalStringFromNull(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

// nullIfEmpty maps empty strings to SQL NULL on write.
func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// nullIfNonPositive maps zero or negative integers to SQL NULL on write.
func nullIfNonPositive(n int) any {
	if n <= 0 {
		return nil
	}
	return n
}
