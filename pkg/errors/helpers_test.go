package errors

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapDatabaseError(t *testing.T) {
	// Test nil error
	assert.Nil(t, WrapDatabaseError(nil))

	// Test sql.ErrNoRows
	err := sql.ErrNoRows
	wrapped := WrapDatabaseError(err)
	assert.NotNil(t, wrapped)
	assert.True(t, Is(wrapped, CodeNotFound))

	// Test duplicate key error
	dupErr := errors.New("duplicate key value violates unique constraint")
	wrapped = WrapDatabaseError(dupErr)
	assert.NotNil(t, wrapped)
	assert.True(t, Is(wrapped, CodeConflict))

	// Test foreign key error
	fkErr := errors.New("violates foreign key constraint")
	wrapped = WrapDatabaseError(fkErr)
	assert.NotNil(t, wrapped)
	assert.True(t, Is(wrapped, CodeInvalidInput))

	// Test connection error
	connErr := errors.New("connection timeout")
	wrapped = WrapDatabaseError(connErr)
	assert.NotNil(t, wrapped)
	assert.True(t, Is(wrapped, CodeDatabaseConnection))

	// Test generic database error
	genericErr := errors.New("some database error")
	wrapped = WrapDatabaseError(genericErr)
	assert.NotNil(t, wrapped)
	assert.True(t, Is(wrapped, CodeDatabaseError))
}

func TestWrapValidationError(t *testing.T) {
	// Test nil error
	assert.Nil(t, WrapValidationError(nil))

	// Test validation error
	err := errors.New("validation failed")
	wrapped := WrapValidationError(err)
	assert.NotNil(t, wrapped)
	assert.True(t, Is(wrapped, CodeValidation))
}

