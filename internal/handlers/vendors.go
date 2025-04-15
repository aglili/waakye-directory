package handlers

import (
	"github.com/aglili/waakye-directory/internal/models"
	"github.com/aglili/waakye-directory/internal/repository/postgres"
	"github.com/aglili/waakye-directory/internal/utils"
	"github.com/gin-gonic/gin"
)

type VendorHandler struct {
	repository        postgres.VendorRepository
	ratingsRepository postgres.RatingsRepository
}

func NewVendorHandler(repository postgres.VendorRepository, ratingsRepository postgres.RatingsRepository) *VendorHandler {
	return &VendorHandler{
		repository:        repository,
		ratingsRepository: ratingsRepository,
	}
}

// CreateVendor godoc
// @Summary Create a new vendor
// @Description Create a new waakye vendor
// @Tags vendors
// @Accept json
// @Produce json
// @Param vendor body CreateWaakyeVendorSchema true "Vendor object"
// @Success 201 {object} CreatedResponse "Vendor created successfully"
// @Failure 400 {object} BadRequestResponse "Bad request"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
// @Router /api/v1/vendors [post]
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
	utils.RespondWithCreated(ctx, createdMessage, vendor)
}

// ListVendorsWithPagination godoc
// @Summary List vendors with pagination
// @Description List all vendors with pagination
// @Tags vendors
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} PaginatedResponse "Vendors retrieved successfully"
// @Failure 400 {object} BadRequestResponse "Bad request"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
// @Router /api/v1/vendors [get]
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
		return
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

// GetVendorByID godoc
// @Summary Get vendor by ID
// @Description Get a vendor by ID
// @Tags vendors
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Success 200 {object} CreatedResponse "Vendor retrieved successfully"
// @Failure 400 {object} BadRequestResponse "Invalid UUID format"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
// @Router /api/v1/vendors/{id} [get]
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


// GetNearbyVendors godoc
// @Summary Get nearby vendors
// @Description Get nearby vendors based on latitude and longitude
// @Tags vendors
// @Accept json
// @Produce json
// @Param lat query float64 true "Latitude"
// @Param lng query float64 true "Longitude"
// @Success 200 {object} CreatedResponse "Nearby vendors retrieved successfully"
// @Failure 400 {object} BadRequestResponse "Invalid latitude or longitude"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
// @Router /api/v1/vendors/nearby [get]
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

// GetVerifiedVendors godoc
// @Summary Get verified vendors
// @Description Get all verified vendors
// @Tags vendors
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} PaginatedResponse "Verified vendors retrieved successfully"
// @Failure 400 {object} BadRequestResponse "Bad request"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
// @Router /api/v1/vendors/verified [get]
func (h *VendorHandler) GetVerifiedVendors(ctx *gin.Context) {
	params, err := utils.GetPaginationParams(ctx)
	if err != nil {
		userMessage := "Failed to list verified vendors"
		utils.RespondWithBadRequest(ctx, err.Error(), userMessage)
		return
	}

	vendors, err := h.repository.GetVerifiedVendors(ctx, params.Page, params.PageSize)
	if err != nil {
		userMessage := "Failed to list verified vendors"
		utils.RespondWithInternalServerError(ctx, err.Error(), userMessage)
		return
	}

	totalItems, err := h.repository.CountVerifiedVendors(ctx)
	if err != nil {
		userMessage := "Failed to list verified vendors"
		utils.RespondWithInternalServerError(ctx, err.Error(), userMessage)
		return
	}

	getMessage := "Verified vendors retrieved successfully"
	utils.SendPaginatedResponse(ctx, vendors, params.Page, params.PageSize, totalItems, getMessage)
}

// RateVendor godoc
// @Summary Rate a vendor
// @Description Rate a vendor by ID
// @Tags vendors
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Param rating body RateVendorRequest true "Rating object"
// @Success 201 {object} CreatedResponse "Vendor rated successfully"
// @Failure 400 {object} BadRequestResponse "Invalid UUID format"
// @Failure 404 {object} NotFoundResponse "Vendor not found"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
// @Router /api/v1/vendors/{id}/rate [post]
func (h *VendorHandler) RateVendor(ctx *gin.Context) {
	parsedUUID, ok := utils.ParseUUID(ctx, "id")
	if !ok {
		return
	}

	var request models.RateVendorRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		userMessage := "Failed to rate vendor"
		utils.RespondWithBadRequest(ctx, err.Error(), userMessage)
		return
	}

	// check if vendor exists
	vendor, err := h.repository.GetVendorByID(ctx, parsedUUID)
	if err != nil {
		userMessage := "Failed to rate vendor"
		utils.RespondWithBadRequest(ctx, err.Error(), userMessage)
		return
	}

	if vendor.ID == (models.WaakyeVendor{}).ID {
		userMessage := "Vendor does not exist"
		utils.RespondWithNotFound(ctx, userMessage, userMessage)
		return
	}

	if err := h.ratingsRepository.RateVendor(ctx,parsedUUID,&request); err != nil {
		userMessage := "Failed to rate vendor"
		utils.RespondWithInternalServerError(ctx, err.Error(), userMessage)
		return
	}

	ratedMessage := "Vendor rated successfully"
	utils.RespondWithCreated(ctx, ratedMessage, request)

}

// GetVendorRatings godoc
// @Summary Get vendor ratings
// @Description Get ratings for a vendor by ID
// @Tags vendors
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Success 200 {object} CreatedResponse "Vendor ratings retrieved successfully"
// @Failure 400 {object} BadRequestResponse "Invalid UUID format"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
// @Router /api/v1/vendors/{id}/ratings [get]
func (h *VendorHandler) GetVendorRatings(ctx *gin.Context) {
	parsedUUID, ok := utils.ParseUUID(ctx, "id")
	if !ok {
		return
	}

	ratings, err := h.ratingsRepository.GetVendorGeneralRatings(ctx, parsedUUID)
	if err != nil {
		userMessage := "Failed to get vendor ratings"
		utils.RespondWithInternalServerError(ctx, err.Error(), userMessage)
		return
	}

	getMessage := "Vendor ratings retrieved successfully"
	utils.RespondWithOK(ctx, getMessage, ratings)
}

// GetTopRatedVendors godoc
// @Summary Get top rated vendors
// @Description Get all top rated vendors
// @Tags vendors
// @Accept json
// @Produce json
// @Success 200 {object} CreatedResponse "Top rated vendors retrieved successfully"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
// @Router /api/v1/vendors/top_rated [get]
func (h *VendorHandler) GetTopRatedVendors(ctx *gin.Context) {

	vendors, err := h.repository.GetTopRatedVendors(ctx)
	if err != nil {
		userMessage := "Failed to get vendor ratings"
		utils.RespondWithInternalServerError(ctx, err.Error(), userMessage)
		return
	}

	getMessage := "Top rated vendors retrieved successfully"
	utils.RespondWithOK(ctx, getMessage, vendors)

}



