package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// AppError represents an application error with code, message, and stack trace
type AppError struct {
	Code       ErrorCode
	Message    string
	HTTPStatus int
	Err        error
	Stack      []Frame
}

// Frame represents a single stack frame
type Frame struct {
	File     string
	Line     int
	Function string
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return string(e.Code)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Format formats the error according to fmt.Formatter
func (e *AppError) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			// Verbose format with stack trace
			fmt.Fprintf(f, "Error Code: %s\n", e.Code)
			fmt.Fprintf(f, "Message: %s\n", e.Message)
			if e.Err != nil {
				fmt.Fprintf(f, "Wrapped Error: %v\n", e.Err)
			}
			fmt.Fprintf(f, "Stack Trace:\n")
			for i, frame := range e.Stack {
				fmt.Fprintf(f, "  [%d] %s\n    %s:%d\n", i, frame.Function, frame.File, frame.Line)
			}
		} else {
			fmt.Fprintf(f, "[%s] %s", e.Code, e.Message)
		}
	case 's':
		fmt.Fprintf(f, e.Message)
	case 'q':
		fmt.Fprintf(f, "%q", e.Message)
	}
}

// captureStack captures the current stack trace, skipping the specified number of frames
func captureStack(skip int) []Frame {
	var frames []Frame
	for i := skip; i < skip+10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}
		frames = append(frames, Frame{
			File:     file,
			Line:     line,
			Function: fn.Name(),
		})
	}
	return frames
}

// New creates a new AppError with the given code and message
func New(code ErrorCode, message string) *AppError {
	httpStatus := GetDefaultRegistry().GetHTTPStatus(code)
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Stack:      captureStack(2),
	}
}

// Newf creates a new AppError with formatted message
func Newf(code ErrorCode, format string, args ...interface{}) *AppError {
	return New(code, fmt.Sprintf(format, args...))
}

// Wrap wraps an existing error with a code and message
func Wrap(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return nil
	}

	httpStatus := GetDefaultRegistry().GetHTTPStatus(code)

	// If err is already an AppError, preserve its stack and wrap it
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			Code:       code,
			Message:    message,
			HTTPStatus: httpStatus,
			Err:        appErr,
			Stack:      captureStack(2),
		}
	}

	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
		Stack:      captureStack(2),
	}
}

// Wrapf wraps an existing error with a formatted message
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *AppError {
	return Wrap(err, code, fmt.Sprintf(format, args...))
}

// WithCode wraps an existing error with just a code (preserves original message)
func WithCode(err error, code ErrorCode) *AppError {
	if err == nil {
		return nil
	}

	httpStatus := GetDefaultRegistry().GetHTTPStatus(code)
	message := err.Error()

	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			Code:       code,
			Message:    message,
			HTTPStatus: httpStatus,
			Err:        appErr,
			Stack:      captureStack(2),
		}
	}

	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
		Stack:      captureStack(2),
	}
}

// Is checks if the error matches the given code
func Is(err error, code ErrorCode) bool {
	if err == nil {
		return false
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}

	return false
}

// GetCode extracts the error code from an error
func GetCode(err error) ErrorCode {
	if err == nil {
		return ""
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}

	return CodeInternalError
}

// GetHTTPStatus extracts the HTTP status code from an error
func GetHTTPStatus(err error) int {
	if err == nil {
		return 500
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.HTTPStatus
	}

	return 500
}

// GetMessage extracts the error message
func GetMessage(err error) string {
	if err == nil {
		return ""
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.Message
	}

	return err.Error()
}

// StackString returns a formatted stack trace string
func (e *AppError) StackString() string {
	var builder strings.Builder
	for i, frame := range e.Stack {
		builder.WriteString(fmt.Sprintf("  [%d] %s\n    %s:%d\n", i, frame.Function, frame.File, frame.Line))
	}
	return builder.String()
}
