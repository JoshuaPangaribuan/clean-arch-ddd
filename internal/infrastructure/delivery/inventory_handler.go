package delivery

import (
	"net/http"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/inventory"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/shared/model"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
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
		appErr := apperrors.New(apperrors.CodeInvalidInput, "Invalid request body: "+err.Error())
		HandleError(c, appErr)
		return
	}

	// Validate input
	if err := h.validator.Struct(input); err != nil {
		HandleValidationError(c, err)
		return
	}

	// Execute use case
	output, err := h.createUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		HandleError(c, err)
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
		appErr := apperrors.New(apperrors.CodeInvalidInput, "Product ID is required")
		HandleError(c, appErr)
		return
	}

	// Execute use case
	output, err := h.getUseCase.Execute(c.Request.Context(), productID)
	if err != nil {
		HandleError(c, err)
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
		appErr := apperrors.New(apperrors.CodeInvalidInput, "Invalid request body: "+err.Error())
		HandleError(c, appErr)
		return
	}

	// Validate input
	if err := h.validator.Struct(input); err != nil {
		HandleValidationError(c, err)
		return
	}

	// Execute use case
	output, err := h.adjustUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, model.NewSuccessResponse(
		"Inventory adjusted successfully",
		output,
	))
}

