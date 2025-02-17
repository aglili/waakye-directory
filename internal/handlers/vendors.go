package handlers

import (
	"strconv"

	"github.com/aglili/waakye-directory/internal/models"
	"github.com/aglili/waakye-directory/internal/repository/postgres"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("Failed to bind JSON")
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.repository.CreateVendor(ctx, &vendor); err != nil {
		log.Error().Err(err).Msg("Failed to create vendor")
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, vendor)
}

func (h *VendorHandler) ListVendorsWithPagination(ctx *gin.Context) {
	// Get the page and pageSize query parameters as strings
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	// Convert page to an integer
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert 'page' to an integer")
		ctx.JSON(400, gin.H{"error": "Invalid 'page' value. Must be an integer."})
		return
	}

	// Convert pageSize to an integer
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert 'page_size' to an integer")
		ctx.JSON(400, gin.H{"error": "Invalid 'page_size' value. Must be an integer."})
		return
	}

	// Ensure page and pageSize are valid
	if page < 1 {
		log.Error().Msg("'page' must be greater than or equal to 1")
		ctx.JSON(400, gin.H{"error": "'page' must be greater than or equal to 1."})
		return
	}
	if pageSize < 1 {
		log.Error().Msg("'page_size' must be greater than or equal to 1")
		ctx.JSON(400, gin.H{"error": "'page_size' must be greater than or equal to 1."})
		return
	}

	// Call the repository method with the validated integers
	vendors, err := h.repository.ListVendorsWithPagination(ctx, page, pageSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list vendors")
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Return the vendors as JSON
	ctx.JSON(200, vendors)
}
