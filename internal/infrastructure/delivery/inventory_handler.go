package delivery

import (
	"errors"
	"net/http"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/inventory"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/shared/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// InventoryHandler handles HTTP requests for inventory operations
type InventoryHandler struct {
	createUseCase *inventory.CreateInventoryUseCase
	getUseCase    *inventory.GetInventoryUseCase
	adjustUseCase *inventory.AdjustInventoryUseCase
	validator     *validator.Validate
}

// NewInventoryHandler creates a new InventoryHandler
func NewInventoryHandler(
	createUseCase *inventory.CreateInventoryUseCase,
	getUseCase *inventory.GetInventoryUseCase,
	adjustUseCase *inventory.AdjustInventoryUseCase,
) *InventoryHandler {
	return &InventoryHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		adjustUseCase: adjustUseCase,
		validator:     validator.New(),
	}
}

// Create handles POST /inventory - creates a new inventory record
func (h *InventoryHandler) Create(c *gin.Context) {
	var input inventory.CreateInventoryInput

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
		// Handle specific errors
		if err.Error() == "cannot create inventory: product not found" {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(
				err.Error(),
				"PRODUCT_NOT_FOUND",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			err.Error(),
			"CREATE_FAILED",
		))
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, model.NewSuccessResponse(
		"Inventory created successfully",
		output,
	))
}

// Get handles GET /inventory/:productId - retrieves inventory by product ID
func (h *InventoryHandler) Get(c *gin.Context) {
	productID := c.Param("productId")

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
		if err.Error() == "inventory not found" {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(
				"Inventory not found for this product",
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
		"Inventory retrieved successfully",
		output,
	))
}

// Adjust handles PATCH /inventory/adjust - adjusts inventory quantity
func (h *InventoryHandler) Adjust(c *gin.Context) {
	var input inventory.AdjustInventoryInput

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
	output, err := h.adjustUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		// Handle specific errors
		if err.Error() == "inventory not found" {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(
				"Inventory not found",
				"NOT_FOUND",
			))
			return
		}
		if err.Error() == "cannot adjust inventory: product not found" {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(
				err.Error(),
				"PRODUCT_NOT_FOUND",
			))
			return
		}
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			err.Error(),
			"ADJUST_FAILED",
		))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, model.NewSuccessResponse(
		"Inventory adjusted successfully",
		output,
	))
}

