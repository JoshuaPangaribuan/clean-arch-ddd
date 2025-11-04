package delivery

import (
	"net/http"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product/command"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product/query"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/shared/model"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ProductHandler handles HTTP requests for product operations
type ProductHandler struct {
	createCommand *command.CreateProductCommand
	getQuery      *query.GetProductQuery
	validator     *validator.Validate
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(
	createCommand *command.CreateProductCommand,
	getQuery *query.GetProductQuery,
) *ProductHandler {
	return &ProductHandler{
		createCommand: createCommand,
		getQuery:      getQuery,
		validator:     validator.New(),
	}
}

// Create handles POST /products - creates a new product
func (h *ProductHandler) Create(c *gin.Context) {
	var input command.CreateProductInput

	// Bind JSON request body
	if err := c.ShouldBindJSON(&input); err != nil {
		appErr := apperrors.New(apperrors.CodeInvalidInput, "Invalid request body: "+err.Error())
		HandleError(c, appErr)
		return
	}

	// Validate input
	if err := h.validator.Struct(input); err != nil {
		HandleValidationError(c, err)
		return
	}

	// Execute command
	output, err := h.createCommand.Execute(c.Request.Context(), input)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, model.NewSuccessResponse(
		"Product created successfully",
		output,
	))
}

// Get handles GET /products/:id - retrieves a product by ID
func (h *ProductHandler) Get(c *gin.Context) {
	productID := c.Param("id")

	if productID == "" {
		appErr := apperrors.New(apperrors.CodeInvalidProductID, "Product ID is required")
		HandleError(c, appErr)
		return
	}

	// Execute query
	output, err := h.getQuery.Execute(c.Request.Context(), productID)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, model.NewSuccessResponse(
		"Product retrieved successfully",
		output,
	))
}

// HealthCheck handles GET /health - simple health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Service is running",
	})
}
