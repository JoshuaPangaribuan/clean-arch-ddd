package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCodeRegistry(t *testing.T) {
	registry := NewErrorCodeRegistry()

	// Test registration
	registry.Register(CodeNotFound, 404, "Not found")
	
	// Test retrieval
	metadata, exists := registry.Get(CodeNotFound)
	assert.True(t, exists)
	assert.Equal(t, CodeNotFound, metadata.Code)
	assert.Equal(t, 404, metadata.HTTPStatus)
	assert.Equal(t, "Not found", metadata.Description)

	// Test non-existent code
	_, exists = registry.Get(ErrorCode("NON_EXISTENT"))
	assert.False(t, exists)

	// Test HTTP status retrieval
	status := registry.GetHTTPStatus(CodeNotFound)
	assert.Equal(t, 404, status)

	// Test default status for unknown code
	status = registry.GetHTTPStatus(ErrorCode("UNKNOWN"))
	assert.Equal(t, 500, status)
}

func TestDefaultRegistry(t *testing.T) {
	registry := GetDefaultRegistry()

	// Test that default codes are registered
	testCases := []struct {
		code       ErrorCode
		httpStatus int
	}{
		{CodeNotFound, 404},
		{CodeInvalidInput, 400},
		{CodeInternalError, 500},
		{CodeConflict, 409},
		{CodeProductNotFound, 404},
		{CodeInventoryNotFound, 404},
	}

	for _, tc := range testCases {
		metadata, exists := registry.Get(tc.code)
		assert.True(t, exists, "Code %s should be registered", tc.code)
		assert.Equal(t, tc.httpStatus, metadata.HTTPStatus, "HTTP status for %s should be %d", tc.code, tc.httpStatus)
	}
}

func TestRegisterErrorCode(t *testing.T) {
	// Register a custom error code
	customCode := ErrorCode("CUSTOM_ERROR")
	RegisterErrorCode(customCode, 418, "I'm a teapot")

	// Verify it's registered
	registry := GetDefaultRegistry()
	metadata, exists := registry.Get(customCode)
	assert.True(t, exists)
	assert.Equal(t, 418, metadata.HTTPStatus)
	assert.Equal(t, "I'm a teapot", metadata.Description)
}

