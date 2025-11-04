package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New(CodeNotFound, "Resource not found")

	assert.NotNil(t, err)
	assert.Equal(t, CodeNotFound, err.Code)
	assert.Equal(t, "Resource not found", err.Message)
	assert.Equal(t, 404, err.HTTPStatus)
	assert.NotEmpty(t, err.Stack)
}

func TestNewf(t *testing.T) {
	err := Newf(CodeInvalidInput, "Invalid input: %s", "test")

	assert.NotNil(t, err)
	assert.Equal(t, CodeInvalidInput, err.Code)
	assert.Equal(t, "Invalid input: test", err.Message)
	assert.Equal(t, 400, err.HTTPStatus)
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrapped := Wrap(originalErr, CodeInternalError, "wrapped error")

	assert.NotNil(t, wrapped)
	assert.Equal(t, CodeInternalError, wrapped.Code)
	assert.Equal(t, "wrapped error", wrapped.Message)
	assert.Equal(t, originalErr, wrapped.Err)
	assert.NotEmpty(t, wrapped.Stack)
}

func TestWrapf(t *testing.T) {
	originalErr := errors.New("original error")
	wrapped := Wrapf(originalErr, CodeInvalidInput, "Invalid: %s", "test")

	assert.NotNil(t, wrapped)
	assert.Equal(t, CodeInvalidInput, wrapped.Code)
	assert.Equal(t, "Invalid: test", wrapped.Message)
	assert.Equal(t, originalErr, wrapped.Err)
}

func TestWithCode(t *testing.T) {
	originalErr := errors.New("original error")
	appErr := WithCode(originalErr, CodeNotFound)

	assert.NotNil(t, appErr)
	assert.Equal(t, CodeNotFound, appErr.Code)
	assert.Equal(t, "original error", appErr.Message)
	assert.Equal(t, originalErr, appErr.Err)
}

func TestWrap_NilError(t *testing.T) {
	wrapped := Wrap(nil, CodeNotFound, "test")
	assert.Nil(t, wrapped)
}

func TestWithCode_NilError(t *testing.T) {
	appErr := WithCode(nil, CodeNotFound)
	assert.Nil(t, appErr)
}

func TestIs(t *testing.T) {
	err := New(CodeNotFound, "Not found")

	assert.True(t, Is(err, CodeNotFound))
	assert.False(t, Is(err, CodeInvalidInput))

	// Test with regular error
	regularErr := errors.New("regular error")
	assert.False(t, Is(regularErr, CodeNotFound))
}

func TestGetCode(t *testing.T) {
	err := New(CodeNotFound, "Not found")
	assert.Equal(t, CodeNotFound, GetCode(err))

	regularErr := errors.New("regular error")
	assert.Equal(t, CodeInternalError, GetCode(regularErr))

	assert.Equal(t, ErrorCode(""), GetCode(nil))
}

func TestGetHTTPStatus(t *testing.T) {
	err := New(CodeNotFound, "Not found")
	assert.Equal(t, 404, GetHTTPStatus(err))

	regularErr := errors.New("regular error")
	assert.Equal(t, 500, GetHTTPStatus(regularErr))

	assert.Equal(t, 500, GetHTTPStatus(nil))
}

func TestGetMessage(t *testing.T) {
	err := New(CodeNotFound, "Not found")
	assert.Equal(t, "Not found", GetMessage(err))

	regularErr := errors.New("regular error")
	assert.Equal(t, "regular error", GetMessage(regularErr))

	assert.Equal(t, "", GetMessage(nil))
}

func TestError_Format(t *testing.T) {
	err := New(CodeNotFound, "Resource not found")

	// Test %s format
	assert.Equal(t, "Resource not found", fmt.Sprintf("%s", err))

	// Test %v format
	assert.Contains(t, fmt.Sprintf("%v", err), string(CodeNotFound))
	assert.Contains(t, fmt.Sprintf("%v", err), "Resource not found")
}

func TestError_Unwrap(t *testing.T) {
	originalErr := errors.New("original")
	wrapped := Wrap(originalErr, CodeInternalError, "wrapped")

	assert.Equal(t, originalErr, wrapped.Unwrap())
}

func TestError_StackString(t *testing.T) {
	err := New(CodeNotFound, "Not found")
	stackStr := err.StackString()

	assert.NotEmpty(t, stackStr)
	assert.Contains(t, stackStr, "errors")
}
