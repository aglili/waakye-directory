package provider

import (
	"database/sql"

	"github.com/aglili/waakye-directory/internal/handlers"
	"github.com/aglili/waakye-directory/internal/repository/postgres"
	cacheRepo "github.com/aglili/waakye-directory/internal/repository/redis"
	"github.com/redis/go-redis/v9"
)

type Provider struct {
	DB            *sql.DB
	Cache         *redis.Client
	VendorHandler *handlers.VendorHandler
}

func NewProvider(db *sql.DB, cache *redis.Client) *Provider {
	vendorRepository := postgres.NewVendorRepository(db)
	vendorCacheRepository := cacheRepo.NewVendorCacheRepository(cache)

	vendorHandler := handlers.NewVendorHandler(vendorRepository, vendorCacheRepository)

	return &Provider{
		DB:            db,
		VendorHandler: vendorHandler,
		Cache:         cache,
	}
}
