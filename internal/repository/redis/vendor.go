package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aglili/waakye-directory/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	vendorKeyPrefix = "vendor:"
	defaultTTL      = 60 * time.Second
)

type VendorCacheRepository interface {
	GetVendorByID(ctx context.Context, id uuid.UUID) (*models.WaakyeVendor, error)
	ExistsInCache(ctx context.Context, cacheKey string) bool
	SetVendorByID(ctx context.Context, id uuid.UUID, vendor *models.WaakyeVendor) error
}

// VendorCacheRepository is a contract for working with vendor data
type vendorCacheRepository struct {
	client *redis.Client
}

func NewVendorCacheRepository(client *redis.Client) VendorCacheRepository {
	return &vendorCacheRepository{
		client: client,
	}
}

// GetVendorByID retrieves a vendor by ID from the cache
func (r *vendorCacheRepository) GetVendorByID(ctx context.Context, id uuid.UUID) (*models.WaakyeVendor, error) {
	cacheKey := vendorKeyPrefix + id.String()

	logger := log.With().
		Str("Operation", "GetVendorByID").
		Str("CacheKey", cacheKey).
		Logger()

	data, err := r.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			logger.Info().Msg("Vendor not found in cache")
			return nil, nil
		}
		logger.Error().Err(err).Msg("Failed to retrieve vendor from cache")
		return nil, err
	}

	var vendor models.WaakyeVendor
	if err := json.Unmarshal(data, &vendor); err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal vendor data")
		return nil, err
	}

	logger.Debug().Msg("Vendor retrieved from cache")
	return &vendor, nil
}

// SetVendorByID stores a vendor in the cache
func (r *vendorCacheRepository) SetVendorByID(ctx context.Context, id uuid.UUID, vendor *models.WaakyeVendor) error {
	cacheKey := vendorKeyPrefix + id.String()

	logger := log.With().
		Str("Operation", "SetVendorByID").
		Str("CacheKey", cacheKey).
		Logger()

	data, err := json.Marshal(vendor)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal vendor data")
		return err
	}

	if err := r.client.Set(ctx, cacheKey, data, defaultTTL).Err(); err != nil {
		logger.Error().Err(err).Msg("Failed to set vendor in cache")
		return err
	}

	logger.Debug().Msg("Vendor cached successfully")
	return nil
}

// ExistsInCache checks if a key exists in the cache
func (r *vendorCacheRepository) ExistsInCache(ctx context.Context, cacheKey string) bool {
	exists, _ := r.client.Exists(ctx, cacheKey).Result()
	return exists > 0
}
