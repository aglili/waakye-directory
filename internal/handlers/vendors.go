package handlers

import (

	"github.com/aglili/waakye-directory/internal/models"
	"github.com/aglili/waakye-directory/internal/repository/postgres"
	"github.com/aglili/waakye-directory/internal/utils"
	"github.com/gin-gonic/gin"
)

type VendorHandler struct {
	repository postgres.VendorRepository
}

func NewVendorHandler(repository postgres.VendorRepository) *VendorHandler {
	return &VendorHandler{
		repository: repository,
	}
}

func (h *VendorHandler) CreateVendor(ctx *gin.Context) {
	var vendor models.WaakyeVendor
	if err := ctx.ShouldBindJSON(&vendor); err != nil {
		userMessage := "Failed to create vendor"
		utils.RespondWithBadRequest(ctx, err.Error(), userMessage)
		return
	}

	if err := h.repository.CreateVendor(ctx, &vendor); err != nil {
		userMessage := "Failed to create vendor"
		utils.RespondWithInternalServerError(ctx, err.Error(), userMessage)
		return
	}

	createdMessage := "Vendor created successfully"
	utils.RespondWithCreated(ctx,createdMessage,vendor)
}

func (h *VendorHandler) ListVendorsWithPagination(ctx *gin.Context) {
	params, err := utils.GetPaginationParams(ctx)
	if err != nil {
		userMessage := "Failed to list vendors with pagination"
		utils.RespondWithBadRequest(ctx, err.Error(), userMessage)
		return
	}

	vendors, err := h.repository.ListVendorsWithPagination(ctx, params.Page, params.PageSize)
	if err != nil {
		userMessage := "Failed to list vendors with pagination"
		utils.RespondWithInternalServerError(ctx, err.Error(), userMessage)
	}

	totalItems, err := h.repository.CountVendors(ctx)
	if err != nil {
		userMessage := "Failed to list vendors with pagination"
		utils.RespondWithInternalServerError(ctx, err.Error(), userMessage)
		return
	}

	getMessage := "Vendors retrieved successfully"
	utils.SendPaginatedResponse(ctx, vendors, params.Page, params.PageSize, totalItems, getMessage)
}

func (h *VendorHandler) GetVendorByID(ctx *gin.Context) {
	parsedUUID, ok := utils.ParseUUID(ctx, "id")
	if !ok {
		return
	}

	vendor, err := h.repository.GetVendorByID(ctx, parsedUUID)
	if err != nil {
		userMessage := "Failed to get vendor"
		utils.RespondWithInternalServerError(ctx, err.Error(), userMessage)
		return
	}

	getMessage := "Vendor retrieved successfully"
	utils.RespondWithOK(ctx, getMessage, vendor)
}



func (h *VendorHandler) GetNearbyVendors(ctx *gin.Context) {
	lat, ok := utils.ParseFloat64(ctx, "lat")
	if !ok {
		return
	}

	lng, ok := utils.ParseFloat64(ctx, "lng")
	if !ok {
		return
	}

	vendors, err := h.repository.GetNearbyVendors(ctx, lat, lng, 5)
	if err != nil {
		userMessage := "Failed to get nearby vendors"
		utils.RespondWithInternalServerError(ctx, err.Error(), userMessage)
		return
	}

	getMessage := "Nearby vendors retrieved successfully"
	utils.RespondWithOK(ctx, getMessage, vendors)
}
