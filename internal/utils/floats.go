package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FloatParsingError represents float parsing error messages
var FloatParsingError = "invalid float value"

// ParseFloat64 parses a float64 from a request parameter
func ParseFloat64(c *gin.Context, param string) (float64, bool) {
	valueStr := c.Param(param)
	if valueStr == "" {
		valueStr = c.Query(param) // Try query parameter if path parameter is empty
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		log.Error().
			Err(err).
			Str("param", param).
			Str("value", valueStr).
			Msg("Failed to parse float64")
		c.JSON(400, gin.H{"error": FloatParsingError})
		return 0, false
	}

	return value, true
}

// ParseFloat64WithRange parses a float64 and validates it's within a specified range
func ParseFloat64WithRange(c *gin.Context, param string, min, max float64) (float64, bool) {
	value, ok := ParseFloat64(c, param)
	if !ok {
		return 0, false
	}

	if value < min || value > max {
		log.Error().
			Float64("value", value).
			Float64("min", min).
			Float64("max", max).
			Str("param", param).
			Msg("Float value out of range")
		c.JSON(400, gin.H{
			"error": "value must be between " + strconv.FormatFloat(min, 'f', -1, 64) +
				" and " + strconv.FormatFloat(max, 'f', -1, 64),
		})
		return 0, false
	}

	return value, true
}

// ParseLatitude parses and validates a latitude value (-90 to 90)
func ParseLatitude(c *gin.Context, param string) (float64, bool) {
	return ParseFloat64WithRange(c, param, -90.0, 90.0)
}

// ParseLongitude parses and validates a longitude value (-180 to 180)
func ParseLongitude(c *gin.Context, param string) (float64, bool) {
	return ParseFloat64WithRange(c, param, -180.0, 180.0)
}
