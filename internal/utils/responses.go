package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// --- Error Response Helpers ---

// RespondWithBadRequest sends a 400 Bad Request response with developer and user messages
func RespondWithBadRequest(ctx *gin.Context, devMessage string, userMessage string) {
	log.Error().Msg(devMessage)
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error":   userMessage,
		"details": devMessage,
	})
}

// RespondWithInternalServerError sends a 500 Internal Server Error response with developer and user messages
func RespondWithInternalServerError(ctx *gin.Context, devMessage string, userMessage string) {
	log.Error().Msg(devMessage)
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error":   userMessage,
		"details": devMessage,
	})
}

// RespondWithCreated sends a 201 Created response with the created resource and a user message
func RespondWithCreated(ctx *gin.Context, userMessage string, data interface{}) {
	ctx.JSON(http.StatusCreated, gin.H{
		"data":    data,
		"message": userMessage,
	})
}

// RespondWithOK sends a 200 OK response with the provided data and a user message
func RespondWithOK(ctx *gin.Context, userMessage string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": userMessage,
	})
}

// RespondWithNotFound sends a 404 Not Found response with developer and user messages
func RespondWithNotFound(ctx *gin.Context, devMessage string, userMessage string) {
	log.Error().Msg(devMessage)
	ctx.JSON(http.StatusNotFound, gin.H{
		"error":   userMessage,
		"details": devMessage,
	})
}

// RespondWithUnauthorized sends a 401 Unauthorized response with developer and user messages
func RespondWithUnauthorized(ctx *gin.Context, devMessage string, userMessage string) {
	log.Error().Msg(devMessage)
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"error":   userMessage,
		"details": devMessage,
	})
}

// --- Pagination Helpers ---

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
	Message    string `json:"message"`
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
		return nil, fmt.Errorf("invalid 'page' value: must be an integer")
	}

	// Convert pageSize to integer
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert 'page_size' to an integer")
		return nil, fmt.Errorf("invalid 'page_size' value: must be an integer")
	}

	// Validate parameters
	if page < 1 {
		log.Error().Msg("'page' must be greater than or equal to 1")
		return nil, fmt.Errorf("'page' must be greater than or equal to 1")
	}
	if pageSize < 1 {
		log.Error().Msg("'page_size' must be greater than or equal to 1")
		return nil, fmt.Errorf("'page_size' must be greater than or equal to 1")
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
func NewPaginatedResponse[T any](data []T, page int, pageSize int, totalItems int64, message string) PaginatedResponse[T] {
	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	return PaginatedResponse[T]{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Message:    message,
	}
}


// SendPaginatedResponse sends a paginated response with appropriate status code
func SendPaginatedResponse[T any](c *gin.Context, data []T, page int, pageSize int, totalItems int64, message string) {
	response := NewPaginatedResponse(data, page, pageSize, totalItems, message)
	c.JSON(http.StatusOK, response) // Directly return the PaginatedResponse
}