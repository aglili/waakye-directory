package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// PaginationParams holds the pagination parameters
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Offset   int `json:"offset"`
}

// PaginatedResponse is a generic struct for paginated data
type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// GetPaginationParams extracts and validates pagination parameters from the request
func GetPaginationParams(c *gin.Context) (*PaginationParams, error) {
	// Get query parameters with defaults
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	// Convert page to integer
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert 'page' to an integer")
		return nil, err
	}

	// Convert pageSize to integer
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert 'page_size' to an integer")
		return nil, err
	}

	// Validate parameters
	if page < 1 {
		log.Error().Msg("'page' must be greater than or equal to 1")
		return nil, err
	}
	if pageSize < 1 {
		log.Error().Msg("'page_size' must be greater than or equal to 1")
		return nil, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	return &PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}, nil
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse[T any](data []T, page int, pageSize int, totalItems int64) PaginatedResponse[T] {
	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	return PaginatedResponse[T]{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

// SendPaginatedResponse sends a paginated response with appropriate status code
func SendPaginatedResponse[T any](c *gin.Context, data []T, page int, pageSize int, totalItems int64) {
	response := NewPaginatedResponse(data, page, pageSize, totalItems)
	c.JSON(200, response)
}
