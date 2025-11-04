package errors

import (
	"database/sql"
	"errors"
)

// WrapDatabaseError wraps database errors with appropriate error codes
func WrapDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	// Check for common database errors
	if errors.Is(err, sql.ErrNoRows) {
		return WithCode(err, CodeNotFound)
	}

	// Check for constraint violations (PostgreSQL specific patterns)
	errStr := err.Error()
	if containsAny(errStr, "duplicate key", "unique constraint", "violates unique constraint") {
		return Wrap(err, CodeConflict, "Resource already exists")
	}

	if containsAny(errStr, "foreign key constraint", "violates foreign key constraint") {
		return Wrap(err, CodeInvalidInput, "Referenced resource does not exist")
	}

	if containsAny(errStr, "connection", "network", "timeout", "dial tcp") {
		return Wrap(err, CodeDatabaseConnection, "Database connection error")
	}

	// Default to generic database error
	return Wrap(err, CodeDatabaseError, "Database operation failed")
}

// WrapValidationError wraps validation errors
func WrapValidationError(err error) error {
	if err == nil {
		return nil
	}
	return Wrap(err, CodeValidation, "Validation failed")
}

// containsAny checks if a string contains any of the given substrings
func containsAny(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

