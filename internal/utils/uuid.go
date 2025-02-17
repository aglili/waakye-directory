package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// ErrInvalidUUID represents an invalid UUID error
var ErrInvalidUUID = "invalid UUID format"

// ParseUUID parses and validates a UUID string from the request
func ParseUUID(c *gin.Context, param string) (uuid.UUID, bool) {
	id := c.Param(param)

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Str("uuid", id).Msg("Failed to parse UUID")
		c.JSON(400, gin.H{"error": ErrInvalidUUID})
		return uuid.Nil, false
	}

	return parsedUUID, true
}

// ValidateUUID checks if a string is a valid UUID without gin context
func ValidateUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

// MustParseUUID parses a UUID string or panics if invalid
// Use this only when you are certain the UUID is valid
func MustParseUUID(id string) uuid.UUID {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		panic("invalid UUID: " + id)
	}
	return parsedUUID
}

// NewUUID generates a new UUID
func NewUUID() uuid.UUID {
	return uuid.New()
}
