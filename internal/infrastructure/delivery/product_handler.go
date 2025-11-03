package delivery

import (
	"errors"
	"net/http"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/shared/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ProductHandler handles HTTP requests for product operations
type ProductHandler struct {
	createUseCase *product.CreateProductUseCase
	getUseCase    *product.GetProductUseCase
	validator     *validator.Validate
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(
	createUseCase *product.CreateProductUseCase,
	getUseCase *product.GetProductUseCase,
) *ProductHandler {
	return &ProductHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		validator:     validator.New(),
	}
}

// Create handles POST /products - creates a new product
func (h *ProductHandler) Create(c *gin.Context) {
	var input product.CreateProductInput

	// Bind JSON request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			"Invalid request body: "+err.Error(),
			"INVALID_REQUEST",
		))
		return
	}

	// Validate input
	if err := h.validator.Struct(input); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			c.JSON(http.StatusBadRequest, model.NewErrorResponse(
				"Validation failed: "+validationErrors.Error(),
				"VALIDATION_ERROR",
			))
			return
		}
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			err.Error(),
			"VALIDATION_ERROR",
		))
		return
	}

	// Execute use case
	output, err := h.createUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			err.Error(),
			"CREATE_FAILED",
		))
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
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			"Product ID is required",
			"INVALID_REQUEST",
		))
		return
	}

	// Execute use case
	output, err := h.getUseCase.Execute(c.Request.Context(), productID)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(
				"Product not found",
				"NOT_FOUND",
			))
			return
		}

		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			err.Error(),
			"GET_FAILED",
		))
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
