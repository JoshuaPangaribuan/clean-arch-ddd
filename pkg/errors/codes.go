package errors

// ErrorCode represents a unique error code identifier
type ErrorCode string

// Predefined error codes
const (
	// Generic errors
	CodeInternalError ErrorCode = "INTERNAL_ERROR"
	CodeInvalidInput  ErrorCode = "INVALID_INPUT"
	CodeNotFound      ErrorCode = "NOT_FOUND"
	CodeConflict      ErrorCode = "CONFLICT"
	CodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	CodeForbidden     ErrorCode = "FORBIDDEN"
	CodeValidation    ErrorCode = "VALIDATION_ERROR"

	// Domain-specific errors - Product
	CodeProductNotFound      ErrorCode = "PRODUCT_NOT_FOUND"
	CodeProductAlreadyExists ErrorCode = "PRODUCT_ALREADY_EXISTS"
	CodeInvalidProductID     ErrorCode = "INVALID_PRODUCT_ID"
	CodeInvalidProductName   ErrorCode = "INVALID_PRODUCT_NAME"
	CodeInvalidPrice         ErrorCode = "INVALID_PRICE"

	// Domain-specific errors - Inventory
	CodeInventoryNotFound ErrorCode = "INVENTORY_NOT_FOUND"
	CodeInventoryExists   ErrorCode = "INVENTORY_ALREADY_EXISTS"
	CodeInsufficientStock ErrorCode = "INSUFFICIENT_STOCK"
	CodeInvalidQuantity   ErrorCode = "INVALID_QUANTITY"
	CodeInvalidAdjustment ErrorCode = "INVALID_ADJUSTMENT"

	// Persistence errors
	CodeDatabaseError      ErrorCode = "DATABASE_ERROR"
	CodeDatabaseConnection ErrorCode = "DATABASE_CONNECTION_ERROR"
	CodeQueryFailed        ErrorCode = "QUERY_FAILED"
	CodeTransactionFailed  ErrorCode = "TRANSACTION_FAILED"
)

// ErrorCodeRegistry holds metadata for error codes
type ErrorCodeRegistry struct {
	codes map[ErrorCode]ErrorCodeMetadata
}

// ErrorCodeMetadata contains metadata for an error code
type ErrorCodeMetadata struct {
	Code        ErrorCode
	HTTPStatus  int
	Description string
}

var globalRegistry *ErrorCodeRegistry

func init() {
	globalRegistry = NewErrorCodeRegistry()
	registerDefaultCodes(globalRegistry)
}

// NewErrorCodeRegistry creates a new error code registry
func NewErrorCodeRegistry() *ErrorCodeRegistry {
	return &ErrorCodeRegistry{
		codes: make(map[ErrorCode]ErrorCodeMetadata),
	}
}

// Register registers a new error code with its metadata
func (r *ErrorCodeRegistry) Register(code ErrorCode, httpStatus int, description string) {
	r.codes[code] = ErrorCodeMetadata{
		Code:        code,
		HTTPStatus:  httpStatus,
		Description: description,
	}
}

// Get retrieves metadata for an error code
func (r *ErrorCodeRegistry) Get(code ErrorCode) (ErrorCodeMetadata, bool) {
	metadata, exists := r.codes[code]
	return metadata, exists
}

// GetHTTPStatus returns the HTTP status code for an error code
func (r *ErrorCodeRegistry) GetHTTPStatus(code ErrorCode) int {
	if metadata, exists := r.codes[code]; exists {
		return metadata.HTTPStatus
	}
	// Default to 500 if code not found
	return 500
}

// GetDefaultRegistry returns the global error code registry
func GetDefaultRegistry() *ErrorCodeRegistry {
	return globalRegistry
}

// RegisterErrorCode registers a new error code in the global registry
// This is the main function to use when adding new error types
func RegisterErrorCode(code ErrorCode, httpStatus int, description string) {
	globalRegistry.Register(code, httpStatus, description)
}

// registerDefaultCodes registers all default error codes
func registerDefaultCodes(registry *ErrorCodeRegistry) {
	// Generic errors
	registry.Register(CodeInternalError, 500, "Internal server error")
	registry.Register(CodeInvalidInput, 400, "Invalid input provided")
	registry.Register(CodeNotFound, 404, "Resource not found")
	registry.Register(CodeConflict, 409, "Resource conflict")
	registry.Register(CodeUnauthorized, 401, "Unauthorized access")
	registry.Register(CodeForbidden, 403, "Forbidden access")
	registry.Register(CodeValidation, 400, "Validation error")

	// Product domain errors
	registry.Register(CodeProductNotFound, 404, "Product not found")
	registry.Register(CodeProductAlreadyExists, 409, "Product already exists")
	registry.Register(CodeInvalidProductID, 400, "Invalid product ID")
	registry.Register(CodeInvalidProductName, 400, "Invalid product name")
	registry.Register(CodeInvalidPrice, 400, "Invalid price")

	// Inventory domain errors
	registry.Register(CodeInventoryNotFound, 404, "Inventory not found")
	registry.Register(CodeInventoryExists, 409, "Inventory already exists")
	registry.Register(CodeInsufficientStock, 400, "Insufficient stock available")
	registry.Register(CodeInvalidQuantity, 400, "Invalid quantity")
	registry.Register(CodeInvalidAdjustment, 400, "Invalid adjustment amount")

	// Persistence errors
	registry.Register(CodeDatabaseError, 500, "Database error")
	registry.Register(CodeDatabaseConnection, 503, "Database connection error")
	registry.Register(CodeQueryFailed, 500, "Query execution failed")
	registry.Register(CodeTransactionFailed, 500, "Transaction failed")
}
