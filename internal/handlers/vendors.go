package handlers

import (
	"github.com/aglili/waakye-directory/internal/models"
	"github.com/aglili/waakye-directory/internal/repository/postgres"
	"github.com/aglili/waakye-directory/internal/utils"
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
	params,err := utils.GetPaginationParams(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get pagination parameters")
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}


	vendors, err := h.repository.ListVendorsWithPagination(ctx, params.Page, params.PageSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list vendors with pagination")
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}


	totalItems, err := h.repository.CountVendors(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count vendors")
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}


	utils.SendPaginatedResponse(ctx, vendors, params.Page, params.PageSize, totalItems)
}
