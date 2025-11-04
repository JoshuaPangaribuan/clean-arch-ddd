package delivery

import (
	"errors"
	"log"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/shared/model"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// HandleError handles errors and returns appropriate HTTP responses
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Extract error code and HTTP status
	code := apperrors.GetCode(err)
	httpStatus := apperrors.GetHTTPStatus(err)
	message := apperrors.GetMessage(err)

	// Log error for debugging (in production, use proper logging)
	log.Printf("Error [%s]: %s", code, message)

	// Return error response
	c.JSON(httpStatus, model.NewErrorResponse(message, string(code)))
}

// HandleValidationError handles validation errors from go-playground/validator
func HandleValidationError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		appErr := apperrors.WrapValidationError(err)
		HandleError(c, appErr)
		return
	}
	// Fallback for non-validation errors
	appErr := apperrors.New(apperrors.CodeValidation, err.Error())
	HandleError(c, appErr)
}

